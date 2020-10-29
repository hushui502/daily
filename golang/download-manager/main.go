package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Download struct {
	Url string
	TargetPath string
	TotalSections int
}

func (d *Download) Do() error {
	fmt.Println("Checking URL")
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Cant process, response is %v", resp.StatusCode))
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-length"))
	if err != nil {
		return err
	}

	var sections = make([][2]int, d.TotalSections)
	eachSize := size / d.TotalSections

	for i := range sections {
		if i == 0 {
			sections[i][0] = 0
		} else {
			sections[i][0] = sections[i-1][1] + 1
		}

		if i < d.TotalSections-1 {
			sections[i][1] = sections[i][0] + eachSize
		} else {
			sections[i][1] = size - 1
		}
	}

	var wg sync.WaitGroup
	for i, s := range sections {
		wg.Add(1)
		go func(i int, s [2]int) {
			defer wg.Done()
			err = d.downloadSection(i, s)
			if err != nil {
				panic(err)
			}
		}(i, s)
	}

	wg.Wait()

	return d.mergeFiles(sections)
}



func (d *Download) downloadSection(i int, c [2]int) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c[0], c[1]))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("cant process, %v", resp.StatusCode))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("section-%v.tmp", i), b, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (d *Download) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		d.Url,
		nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", "Silly Download Manager v001")
	return r, nil
}

func (d *Download) mergeFiles(sections [][2]int) error {
	f, err := os.OpenFile(d.TargetPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {

	}

	defer f.Close()

	for i := range sections {
		tmpFileName := fmt.Sprintf("section-%v.tmp", i)
		b, err := ioutil.ReadFile(tmpFileName)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		os.Remove(tmpFileName)
	}

	return nil
}

func main() {
	//startTime := time.Now()
	d := Download{
		Url: "https://sourceforge.net/projects/mingw-w64/files/latest/download",
		TargetPath: "D:\\project\\go\\src\\awesomeProject2\\download-manager",
		TotalSections: 10,
	}

	err := d.Do()
	if err != nil {
		log.Println()
	}
}
