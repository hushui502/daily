package main

import (
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"os"
	"sync"
	"time"
)

//func main() {
//	e := email.NewEmail()
//	e.From = "hf <17862972248@163.com>"
//	e.To = []string{"502131010@qq.com"}
//	e.Subject = "aws"
//	e.HTML = []byte(`
//  <ul>
//<li><a "https://darjun.github.io/2020/01/10/godailylib/flag/">Go 每日一库之 flag</a></li>
//<li><a "https://darjun.github.io/2020/01/10/godailylib/go-flags/">Go 每日一库之 go-flags</a></li>
//<li><a "https://darjun.github.io/2020/01/14/godailylib/go-homedir/">Go 每日一库之 go-homedir</a></li>
//<li><a "https://darjun.github.io/2020/01/15/godailylib/go-ini/">Go 每日一库之 go-ini</a></li>
//<li><a "https://darjun.github.io/2020/01/17/godailylib/cobra/">Go 每日一库之 cobra</a></li>
//<li><a "https://darjun.github.io/2020/01/18/godailylib/viper/">Go 每日一库之 viper</a></li>
//<li><a "https://darjun.github.io/2020/01/19/godailylib/fsnotify/">Go 每日一库之 fsnotify</a></li>
//<li><a "https://darjun.github.io/2020/01/20/godailylib/cast/">Go 每日一库之 cast</a></li>
//</ul>
//  `)
//	err := e.Send("smtp.163.com:25", smtp.PlainAuth("", "17862972248@163.com", "hufan123", "smtp.163.com"))
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	ch := make(chan *email.Email, 10)
	p, err := email.NewPool(
		"smtp.163.com:25",
		4,
		smtp.PlainAuth("", "17862972248@163.com", "hufan123", "smtp.163.com"),
		)
	if err != nil {
		log.Fatal("failed to create pool")
	}

	var wg sync.WaitGroup
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for e := range ch {
				err := p.Send(e, time.Second*10)
				if err != nil {
					fmt.Fprintf(os.Stderr, "email:%v sent error:%v\n", e, err)
				}
			}
		}()
	}

	for i := 0; i < 50; i++ {
			e := email.NewEmail()
			e.From = "hf <17862972248@163.com>"
			e.To = []string{"502131010@qq.com"}
			e.Subject = "aws"
			e.Text = []byte("test pool")
			ch <- e
	}

	close(ch)
	wg.Wait()
}