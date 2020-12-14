package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

func lineByLine(file string) error {
	var err error

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("read file error %s", err)
			break
		}
		fmt.Println(line)
	}

	return nil
}

func wordByWord(file string) error {
	var err error

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("read file error %s", err)
			break
		}
		r := regexp.MustCompile("[^\\s]+")
		words := r.FindAllString(line, -1)
		for i := 0; i < len(words); i++ {
			fmt.Printf(words[i])
		}
	}

	return nil
}

func charByChar(file string) error {
	var err error
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			return err
		}
		for _, x := range line {
			fmt.Println(string(x))
		}
	}
	return nil
}

func readSize(f *os.File, size int) []byte {
	buffer := make([]byte, size)

	n, err := f.Read(buffer)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return buffer[:n]
}

func writeFile() {
	s := []byte("data 1-2-3")
	f, err := os.Create("f1.txt")
	if err != nil {
		return
	}
	defer f.Close()

	// 1
	fmt.Fprintf(f, string(s))

	// 2
	w := bufio.NewWriter(f)
	w.WriteString(string(s))

	// 3
	ioutil.WriteFile("f1.txt", s, 0644)

	// 4
	io.WriteString(f, string(s))
}

func save(fileName string, data []byte) error {
	err := os.Remove(fileName)
	if err != nil {
		return err
	}

	saveTo, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer saveTo.Close()

	encoder := gob.NewEncoder(saveTo)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func load() error {
	res := []byte{}

	loadFrom, err := os.Open("")
	defer loadFrom.Close()

	decoder := gob.NewDecoder(loadFrom)
	decoder.Decode(&res)

	return nil
}



func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Printf("usage: byLine <file1> [<file2> ...]\n")
		return
	}
	for _, file := range flag.Args() {
		err := lineByLine(file)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("--------")
	arguments := os.Args
	if len(arguments) != 3 {
		fmt.Println("<buffer size> <filename>")
		return
	}
	bufferSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Open(os.Args[2])
	defer file.Close()

	for {
		readData := readSize(f, bufferSize)
		if readData != nil {
			fmt.Println(string(readData))
		}
		break
	}

}


