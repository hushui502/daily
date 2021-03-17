package main

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func initCmd(cmd *cobra.Command, args []string) {
	s := bufio.NewScanner(os.Stdin)
	log.Printf("SMTP Host Address: ")
	if !s.Scan() {
		log.Println("init was canceled!")
		return
	}
	info := s.Text()
	tcli.SMTPHost = info

	log.Printf("SMTP Host Port: ")
	if !s.Scan() {
		log.Println("init was canceled.")
		return
	}
	info = s.Text()
	tcli.SMTPPort = info

	log.Printf("Avatar (Your Name): ")
	if !s.Scan() {
		log.Println("init was canceled.")
		return
	}
	info = s.Text()
	tcli.Avatar = info

	log.Printf("Email (Your Email): ")
	if !s.Scan() {
		log.Println("init was canceled.")
		return
	}
	info = s.Text()
	tcli.EmailAddr = info

	log.Printf("Username (Your Email's Username): ")
	if !s.Scan() {
		log.Println("init was canceled.")
		return
	}
	info = s.Text()
	tcli.Username = info

	log.Printf("Password (Your Email's Password): ")
	if !s.Scan() {
		log.Println("init was canceled.")
		return
	}
	info = s.Text()
	tcli.Password = info

	log.Printf("Things 3 Email Address: ")
	if !s.Scan() {
		log.Println("init was canceled.")
		return
	}
	info = s.Text()
	tcli.ThingsAddr = info
	tcli.save()
	log.Println("You can start using tli now :)")
}

func logCmd(cmd *cobra.Command, args []string) {
	checkhome()
	tcli.parse()

	var (
		n int
		err error
	)

	if len(args) != 0 {
		n, err = strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("invalid argument, please input a number.")
		}
	}

	data, err := ioutil.ReadFile(homedir + "/" + pathHist)
	if err != nil {
		log.Fatalf("cannot read ~/.tli_history, try store something first.")
	}

	d := yaml.NewDecoder(bytes.NewReader(data))

	rs := []record{}
	for {
		var r record
		err = d.Decode(&r)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("corrupted ~/.tli_history file, err: %v", err)
		}
		rs = append(rs, r)
	}

	if n == 0 {
		n = len(rs)
	}
	for i := 1; i <= n; i++ {
		out, _ := yaml.Marshal(rs[len(rs)-i])
		log.Println(string(out))
	}
}

func todoCmd(cmd *cobra.Command, args []string) {
	checkhome()
	tcli.parse()

	title := strings.Join(args, " ")
	a, err := newtcliTODO(title)
	if errors.Is(err, errCanceled) {
		log.Println("TODO is canceled.")
		return
	}
	a.Save()
	a.Range(func(title, body string) {
		for i := 0; i < 5; i++ {
			err := tcli.sendIndex(title, body)
			if err == nil {
				return
			}
			log.Printf("failed to send inbox, err: %v. Retry...", err)
		}
	})

	log.Println("Done!")
}