package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mime"
	"net/smtp"
	"os"
	"os/signal"
	"os/user"
	"runtime"
	"strings"
	"time"
)

// build info
var (
	Version     string
	BuildTime   string
	GoVersion   = runtime.Version()
	errCanceled = errors.New("action canceled")
	homedir     string
	pathConf    = ".tcli_config"
	pathHist    = ".tcli_history"
)

var tcli tcliconf

func checkhome() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(`can not find home directory, err: %v`, err)
	}
	homedir = user.HomeDir
	if len(homedir) == 0 {
		log.Fatalf(`can not find home directory, err: %v`, err)
	}
}

// tcliconf contains all necessary information for sending a thing
// to the thing inbox
type tcliconf struct {
	SMTPHost   string `yaml:"smtp_host"`
	SMTPPort   string `yaml:"smtp_port"`
	Avatar     string `yaml:"avatar"`
	EmailAddr  string `yaml:"email_addr"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	ThingsAddr string `yaml:"things_addr"`
}

func (c *tcliconf) parse() {
	f := os.Getenv("TLE_CONF")
	d, err := ioutil.ReadFile(f)
	if err != nil {
		d, err = ioutil.ReadFile(homedir + "/" + pathConf)
		if err != nil {
			log.Fatalf(`cannot find tli config, err: %v try: ctli init`, err)
		}
	}
	err = yaml.Unmarshal(d, c)
	if err != nil {
		log.Fatalf(`cannot parse tli config, err: %v`, err)
	}
}

func (c *tcliconf) save() {
	checkhome()
	data, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("cannot save your data, err: %v", err)
	}

	f, err := os.OpenFile(homedir+"/"+pathConf,
		os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	all := []byte("---\n")
	all = append(all, data...)
	if _, err := f.Write(all); err != nil {
		return
	}
}

func (c *tcliconf) sendIndex(title, body string) error {
	if strings.ContainsAny(title, "\"#$%&'(),.:;<>@[]^`{|}~") {
		title = mime.BEncoding.Encode("utf-8", title)
	} else {
		title = mime.QEncoding.Encode("utf-8", title)
	}

	err := smtp.SendMail(
		c.SMTPHost + ":" + c.SMTPPort,
		smtp.PlainAuth("", c.Username, c.Password, c.SMTPHost),
		c.EmailAddr, []string{c.ThingsAddr},
		[]byte(fmt.Sprintf("Subject: %s\r\nFrom: %s <%s>\r\nTo: %s\r\n%s",
			title,
			c.Avatar,
			c.EmailAddr,
			c.ThingsAddr,
			body,
		)))
	if err != nil {
		return err
	}

	return nil
}

type tcliTODO struct {
	title string
	body []string
}

func newtcliTODO(title string) (*tcliTODO, error) {
	a := &tcliTODO{
		title: title,
	}
	if !a.waitBody() {
		return nil, errCanceled
	}

	return a, nil
}

func (a *tcliTODO) waitBody() bool {
	s := bufio.NewScanner(os.Stdin)
	fmt.Println("(Enter an empty line to complete; Ctrl+C/Ctrl+D to cancel)")

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, os.Interrupt)

	line := make(chan string, 1)
	go func() {
		for {
			fmt.Println("> ")
			if !s.Scan() {
				signCh <- os.Interrupt
				return
			}
			l := s.Text()
			if len(l) == 0 {
				line <- ""
				return
			}
			line <- l
		}
	}()

	for {
		select {
		case <-signCh:
			return false
		case l := <-line:
			if len(l) == 0 {
				return true
			}
			a.body = append(a.body, l)
		}
	}
}

const maxlen = 2000

func (a *tcliTODO) Range(f func(string, string)) {
	whole := strings.Join(a.body, "\n")

	if len(whole) < maxlen {
		f(a.title, whole)
	}

	count := 1
	for i := 0; i < len(whole); i += maxlen  {
		f(a.title+fmt.Sprintf(" (%d)", count), whole[i:min(i+maxlen, len(whole))])
		count++
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


type record struct {
	Time  time.Time `yaml:"time"`
	Title string    `yaml:"title"`
	Body  string    `yaml:"body"`
}

func (a *tcliTODO) Save() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("cannot save your TODO, err: %v", err)
		}
	}()

	r := record{
		Time: time.Now(),
		Title: a.title,
		Body: strings.Join(a.body, "\n"),
	}

	data, err := yaml.Marshal(r)
	if err != nil {
		return
	}

	f, err := os.OpenFile(homedir+"/"+pathHist,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	all := []byte("---\n")
	all = append(all, data...)
	all = append(all, []byte("\n")...)
	if _, err := f.Write(all); err != nil {
		return
	}
}

func main() {
	log.SetPrefix("tcli: ")
	log.SetFlags(log.Ldate | log.Ltime)

	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "initialize tli settings",
		Long:  `init will ask you several informations for setting up the configuration.`,
		Args:  cobra.ExactArgs(0),
		Run:   initCmd,
	}


	var cmdLog = &cobra.Command{
		Use:   "log [number]",
		Short: "print logs",
		Long:  `log will print the specified number of items`,
		Args:  cobra.MinimumNArgs(0),
		Run:   logCmd,
	}

	var cmdTodo = &cobra.Command{
		Use:                   "todo [title]",
		Short:                 "create a todo and send it to the Things' Inbox",
		Long:                  "create a todo and send it to the Things' Inbox.",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run:                   todoCmd,
	}

	var rootCmd = &cobra.Command{
		Use:   "tli",
		Short: "A Things CLI for Linux support.",
		Long: fmt.Sprintf(`
tli is a Linux CLI that supports send items to the Things' Inbox safely.
Specifically, it will save the sent TODO log to prevent if you send too
much to the Things' server. tli also checks your content to make sure your
inputs won't be too large so that the content is not silently truncated
by Things.
BuildVersion: %v
BuildTime:    %v
GoVersion:    %v
`, Version, BuildTime, GoVersion),
	}


	rootCmd.AddCommand(cmdInit)
	rootCmd.AddCommand(cmdLog)
	rootCmd.AddCommand(cmdTodo)
	rootCmd.Execute()
}
