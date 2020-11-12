package main

import (
	"fmt"
	"lib/simplepow/pow"
	"lib/simplepow/util"
	"log"
	"math/big"
)
var logs = log.Logger{}

func main() {
	forever := make(chan string)

	powManager, err := pow.NewProofOfWorkManager(logs, pow.WithThreads(5))
	if err != nil {
		panic(err)
	}
	result := make(chan *pow.Result)
	err = powManager.Work(&pow.Block{
		Header:&pow.Header{
			Difficulty: new(big.Int).Exp(big.NewInt(2), big.NewInt(8 * 3), big.NewInt(0)),
		},
	}, result)
	if err != nil {
		panic(err)
	}
	select {
	case result_ := <-result:
		fmt.Println(result_.Nonce, result_.AttemptNum, util.BufferToHexString(result_.Hash, true))
	}

	<-forever
}
