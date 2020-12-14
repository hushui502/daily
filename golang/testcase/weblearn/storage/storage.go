package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

func main() {
	csvFile, err := os.Create("posts.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	allPosts := []Post{
		{Id: 1, Content: "hello", Author: "hufan"},
		{Id: 2, Content: "world", Author: "libai"},
	}

	writer := csv.NewWriter(csvFile)
	for _, post := range allPosts {
		line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
		err := writer.Write(line)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()

	file, err := os.Open("posts.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var posts []Post
	for _, item := range record {
		id, _ := strconv.Atoi(item[0])
		//id, _ := strconv.ParseInt(item[0], 0,0)
		post := Post{id, item[1], item[2]}
		posts = append(posts, post)
	}

	fmt.Println(posts)

	//a1 := Post{Id:12, Content:"dsd", Author:"sds"}
	Store(posts, "posts.temp")
	var pp []Post
	Load(&pp, "posts.temp")
	fmt.Println(pp)
}

func Store(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func Load(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}
