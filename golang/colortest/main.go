package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var (
	pass = color.FgGreen
	skip = color.FgYellow
	fail = color.FgRed

	skipnotest bool
)

const (
	paletteEnv = "GOTEST_PALETTE"
	skipNoTestsEnv = "GOTEST_SKIPNOTESTS"
)

func main() {
	enablePalette()
	enableSkipNoTests()
	enableOnCI()

	// skip *.go
	os.Exit(gotest(os.Args[1:]))
}

func gotest(args []string) int {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()

	// pipe instead of temp buffer
	r, w := io.Pipe()
	// only w is closed ==> return
	defer w.Close()

	args = append([]string{"test"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = w
	cmd.Stdout = w
	// set env
	cmd.Env = os.Environ()

	go consume(&wg, r)

	sigc := make(chan os.Signal)
	done := make(chan struct{})
	// defer close loop
	defer func() {
		done <- struct{}{}
	}()
	signal.Notify(sigc)

	// monitor for this cmd exec
	go func() {
		for {
			select {
			case sig := <-sigc:
				cmd.Process.Signal(sig)
			case <-done:
				return
			}
		}
	}()

	if err := cmd.Run(); err != nil {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			return ws.ExitStatus()
		}
		return 1
	}

	return 0
}

func consume(wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()
	reader := bufio.NewReader(r)
	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Print(err)
			return
		}
		parse(string(l))
	}
}

func parse(line string) {
	trimmed := strings.TrimSpace(line)
	defer color.Unset()

	var c color.Attribute
	switch {
	case strings.Contains(trimmed, "[no test files]"):
		if skipnotest {
			return
		}
	case strings.HasPrefix(trimmed, "--- PASS"):	// passed
		fallthrough
	case strings.HasPrefix(trimmed, "ok"):
		fallthrough
	case strings.HasPrefix(trimmed, "PASS"):
		c = pass
	case strings.HasPrefix(trimmed, "--- SKIP"):	// skipped
		c = skip
	case strings.HasPrefix(trimmed, "FAIL"):	// failed
		c = fail
	}
	color.Set(c)
	fmt.Printf("%s\n", line)
}

func enableOnCI() {
	ci := strings.ToLower(os.Getenv("CI"))
	switch ci {
	case "true":
		fallthrough
	case "travis":
		fallthrough
	case "appveyor":
		fallthrough
	case "circlecli":
		color.NoColor = false
	}
}

func enablePalette() {
	v := os.Getenv(paletteEnv)
	if v == "" {
		return
	}
	vals := strings.Split(v, ",")
	if len(vals) != 2 {
		return
	}
	// todo
	if c, ok := colors[vals[0]]; ok {
		fail = c
	}
	if c, ok := colors[vals[1]]; ok {
		pass = c
	}
}

func enableSkipNoTests() {
	v := os.Getenv(skipNoTestsEnv)
	if v == "" {
		return
	}
	v = strings.ToLower(v)
	skipnotest = v == "true"
}

var colors = map[string]color.Attribute{
	"black":     color.FgBlack,
	"hiblack":   color.FgHiBlack,
	"red":       color.FgRed,
	"hired":     color.FgHiRed,
	"green":     color.FgGreen,
	"higreen":   color.FgHiGreen,
	"yellow":    color.FgYellow,
	"hiyellow":  color.FgHiYellow,
	"blue":      color.FgBlue,
	"hiblue":    color.FgHiBlue,
	"magenta":   color.FgMagenta,
	"himagenta": color.FgHiMagenta,
	"cyan":      color.FgCyan,
	"hicyan":    color.FgHiCyan,
	"white":     color.FgWhite,
	"hiwhite":   color.FgHiWhite,
}