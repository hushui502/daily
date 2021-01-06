package main

import (
	"bytes"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	mdata, err := asn1.Marshal(12)
	if err != nil {
		return
	}

	var n int
	_, err = asn1.Unmarshal(mdata, &n)
	println("After marshal/unmarshal: ", n)


	s := "hello"
	mstring, _ := asn1.Marshal(s)

	var newstr string
	_, err = asn1.Unmarshal(mstring, &newstr)
	println("After marshal/unmarshal: ", newstr)
}

func encodeNetConn()  {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		daytime := time.Now()
		mdata, _ := asn1.Marshal(daytime)
		conn.Write(mdata)
		conn.Close()
	}
}

func decodeNetConn() {
	// ":1200"
	service := os.Args[1]
	conn, err := net.Dial("tcp", service)
	if err != nil {
		return
	}

	result, err := readFully(conn)
	if err != nil {
		return
	}

	var newtime time.Time
	_, err = asn1.Unmarshal(result, &newtime)
	if err != nil {
		return
	}

	fmt.Println("After marshal/unmarshal is ", newtime.String())

	os.Exit(1)
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return result.Bytes(), nil
}

// --------------------JSON------------------------
type Person struct {
	Name Name
	Email []Email
}

type Name struct {
	Family string
	Personal string
}

type Email struct {
	Kind string
	Address string
}

func (p Person) String() string {
	s := p.Name.Personal + " " + p.Name.Family
	for _, v := range p.Email {
		s += "\n" + v.Kind + ": " + v.Address
	}

	return s
}

func encodeJson()  {
	person := Person{
		Name:  Name{Family: "f1", Personal: "jan"},
		Email: []Email{
			{Kind: "home", Address: "jan@home.com"},
			{Kind: "work", Address: "jan@work.com"},
		},
	}
	saveJson("person.json", person)
}

func decodeJson() {
	var person Person
	loadJson("person.json", &person)

	// customer string
	fmt.Println(person.String())
}

func loadJson(filename string, key interface{})  {
	inFile, err := os.Open(filename)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(key)
	if err != nil {
		return
	}
	inFile.Close()
}

func saveJson(filename string, key interface{}) {


	outFile, err := os.Create(filename)
	if err != nil {
		return
	}
	encoder := json.NewEncoder(outFile)
	err = encoder.Encode(key)
	outFile.Close()
}