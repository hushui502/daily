package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func parseFileInfoFrom(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			panic(err)
		}

		return params["filename"]
	}
	filename := filepath.Base(resp.Request.URL.Path)

	return filename
}

type FileDownloader struct {
	fileSize int
	url string
	outputFileName string
	totalPart int
	outputDir string
	doneFilePart []filePart
}

func NewFileDownloader(url, outputFileName, outputDir string, totalPart int) *FileDownloader {
	if outputDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		outputDir = wd
	}

	return &FileDownloader{ 
		fileSize:       0,
		url:            url,
		outputFileName: outputFileName,
		outputDir:      outputDir,
		totalPart:      totalPart,
		doneFilePart:   make([]filePart, totalPart),
	}
}

type filePart struct {
	Index int
	From int
	To int
	Data []byte
}

func (d *FileDownloader) head() (int, error) {
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode > 299 {
		return 0, errors.New("")
	}
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return 0, err
	}
	d.outputFileName = parseFileInfoFrom(resp)
	return strconv.Atoi(resp.Header.Get("Content-length"))
}

func (d *FileDownloader) Run() error {
	fileTotalSize, err := d.head()
	if err != nil {
		return err
	}
	d.fileSize = fileTotalSize

	jobs := make([]filePart, d.totalPart)
	eachSize := fileTotalSize / d.totalPart

	for i := range jobs {
		jobs[i].Index = i
		if i == 0 {
			jobs[i].From = 0
		} else {
			jobs[i].From = jobs[i-1].To - 1
		}
		if i < d.totalPart-1 {
			jobs[i].To = jobs[i].From + eachSize
		} else {
			jobs[i].To = fileTotalSize - 1
		}
	}

	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(job filePart) {
			defer wg.Done()
			err := d.downloadPart(job)
			if err != nil {
				log.Println("download failed")
			}
		}(j)
	}

	wg.Wait()
	return d.mergeFileParts()
}

func (d *FileDownloader) downloadPart(c filePart) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c.From, c.To))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("服务器错误状态码: %v", resp.StatusCode))
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(bs) != (c.To - c.From + 1) {
		return errors.New("download file splice len err")
	}

	c.Data = bs
	d.doneFilePart[c.Index] = c
	return nil
}

func (d *FileDownloader) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		d.url,
		nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "mojocn")
	return r, nil
}

func (d *FileDownloader) mergeFileParts() error {
	log.Println("start merging")
	path := filepath.Join(d.outputDir, d.outputFileName)
	mergeFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer mergeFile.Close()

	hash := sha256.New()
	totalSize := 0
	for _, s := range d.doneFilePart {
		mergeFile.Write(s.Data)
		hash.Write(s.Data)
		totalSize += len(s.Data)
	}

	if totalSize != d.fileSize {
		return errors.New("file is not complete")
	}
	if hex.EncodeToString(hash.Sum(nil)) != "" {
		return errors.New("file is bad")
	}

	return nil
}

func main() {
	startTime := time.Now()
	var url string
	url = ""
	downloader := NewFileDownloader(url, "", "", 10)
	if err := downloader.Run(); err != nil {

	}
	fmt.Printf("total spend time %f second", time.Now().Sub(startTime).Seconds())
}
