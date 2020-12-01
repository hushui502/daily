package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// clone this repo
// navigate to this directory and rename the example file mv example.env .env
// go run main.go
// open a new terminal window and nc localhost 9000
// input a token amount to set stake
// input a BPM
// wait a few seconds to see which of the two terminals won
// open as many terminal windows as you like and nc localhost 9000 and watch Proof of Stake in action!

type Block struct {
	Index int
	TimeStamp string
	BPM int
	Hash string
	PrevHash string
	// 和一个global map绑定
	Validator string
}

// Blockchain是我们的正式的区块链，这只是一系列经过验证的块
var Blockchain []Block

// tempBlocks仅仅是一个区块的存储箱，然后在其中一个被选为获胜者将被添加到Blockchain
var tempBlocks []Block

// candidateBlocks是一个块的通道；每个节点提出一个新的块将它发送到这个通道
var candidateBlocks = make(chan Block)

// announcements是一个通道，在这里，我们主要的TCP服务器向所有节点广播最新的BooStand链
var announcements = make(chan string)

var mutex = &sync.Mutex{}

// validators是节点的映射和它们所持有的令牌的数量
var validators = make(map[string]int)

func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.TimeStamp + string(block.BPM) + block.PrevHash

	return calculateHash(record)
}

func generateBlock(oldBlock Block, BPM int, address string) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.TimeStamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = address

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()

	// validate address
	var address string

	io.WriteString(conn, "Enter token balance: ")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = calculateHash(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	io.WriteString(conn, "\nEnter a new BPM: ")

	scanBPM := bufio.NewScanner(conn)

	go func() {
		for {
			for scanBPM.Scan() {
				bpm, err := strconv.Atoi(scanBPM.Text())
				if err != nil {
					log.Printf("%v not a number: %v", scanBPM.Text(), err)
					// bad node, so delete its validate
					delete(validators, address)
					conn.Close()
				}

				mutex.Lock()
				oldLastBlock := Blockchain[len(Blockchain)-1]
				mutex.Unlock()

				newBlock, err := generateBlock(oldLastBlock, bpm, address)
				if err != nil {
					log.Println(err)
					continue
				}
				if isBlockValid(newBlock, oldLastBlock) {
					candidateBlocks <- newBlock
				}

				io.WriteString(conn, "\nEnter a new BPM: ")
			}
		}
	}()

	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(conn, string(output)+"\n")
	}
}

func pickWinner() {
	time.Sleep(30*time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp) > 0 {
		OUTER:
			for _, block := range temp {
				for _, node := range lotteryPool {
					if block.Validator == node {
						continue OUTER
					}
				}

				mutex.Lock()
				setValidators := validators
				mutex.Unlock()

				k, ok := setValidators[block.Validator]
				if ok {
					for i := 0; i < k; i++ {
						lotteryPool = append(lotteryPool, block.Validator)
					}
				}
			}

			s := rand.NewSource(time.Now().Unix())
			r := rand.New(s)
			lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

			for _, block := range temp {
				if block.Validator == lotteryWinner {
					mutex.Lock()
					Blockchain = append(Blockchain, block)
					mutex.Unlock()

					for _ = range validators {
						announcements <- "\nwinning validator: " + lotteryWinner + "\n"
					}
					break
				}
			}
	}

	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateBlockHash(genesisBlock), "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	tcpPort := os.Getenv("PORT")

	server, err := net.Listen("tcp", ":" + tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP Server Listening on port :", tcpPort)
	defer server.Close()

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			pickWinner()
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConn(conn)
	}

}


