package main

import (
	"encoding/gob"
	"os"
)

func saveGob(filename string, key interface{}) {
	outFile, err := os.Create(filename)
	if err != nil {
		return
	}
	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	outFile.Close()
}

func loadGob(filename string, key interface{}) {
	inFile, err := os.Open(filename)
	if err != nil {
		return
	}
	decoder := gob.NewDecoder(inFile)
	err = decoder.Decode(key)
	inFile.Close()
}


