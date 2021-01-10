package ffpprof

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

// Profile represents a pprof profile
type Profile interface {
	Capture() (profile string, err error)
}

// CPUProfile captures the CPU profile
type CPUProfile struct {
	Duration time.Duration	// 30 seconds by default
}

func (c CPUProfile) Capture() (string, error) {
	dur := c.Duration
	if dur == 0 {
		dur = 30 * time.Second
	}

	// create a tmp file for save profile data
	f := newTemp()
	if err := pprof.StartCPUProfile(f); err != nil {
		return "", nil
	}
	// auto stop capture in (duration + current-time)
	time.Sleep(dur)
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		return "", nil
	}

	return f.Name(), nil
}

type HeapProfile struct {}

func (p HeapProfile) Capture() (string, error) {
	return captureProfile("heap")
}

type MutexProfile struct {}

func (p MutexProfile) Capture() (string, error) {
	return captureProfile("mutex")
}

type BlockProfile struct {
	Rate int
}

func (p BlockProfile) Capture() (string, error) {
	if p.Rate > 0 {
		runtime.SetBlockProfileRate(p.Rate)
	}
	return captureProfile("block")
}

// goroutine profile captures the stack traces of all current goroutines
type GoroutineProfile struct {}

func (p GoroutineProfile) Capture() (string, error) {
	return captureProfile("goroutine")
}

// threadcreate profile captures the stack traces that let to the creation of new os threads
type ThreadCreateProfile struct{}

func (p ThreadCreateProfile) Capture() (string, error) {
	return captureProfile("threadcreate")
}

func captureProfile(name string) (string, error) {
	f := newTemp()
	if err := pprof.Lookup(name).WriteTo(f, 2); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	return f.Name(), nil
}

// Capture captures the given profiles at SIGINT
// and opens a browser with the collected profiles.
//
// Capture should be used in development-time
// and should not be used in production binaries.
func Capture(p Profile) {
	go capture(p)
}

func capture(p Profile) {
	// avoid of blocking
	c := make(chan os.Signal, 1)

	switch runtime.GOOS {
	case "windows":
		signal.Notify(c, os.Interrupt)
		fmt.Println("Send interrupt (CTRL-BREAK) to the process to capture")
		fmt.Printf("Use (taskkill /F /PID %d) to end process\n", os.Getpid())
	default:
		signal.Notify(c, syscall.SIGQUIT)
		fmt.Println("Send SIGQUIT (CTRL+\\) to the process to capture...")
	}

	for {
		<-c
		log.Println("Starting to capture.")

		profile, err := p.Capture()
		if err != nil {
			log.Printf("Cannot capture profile: %v", err)
		}

		// open profile with pprof tool
		log.Printf("Starting go tool pprof %v", profile)
		cmd := exec.Command("go", "tool", "pprof", "-http=:", profile)
		if err := cmd.Run(); err != nil {
			log.Printf("Cannot start pprof web UI: %v", err)
		}
	}
}

// Create a temp profile file and return file-handler
func newTemp() (f *os.File) {
	// https://golang.org/src/io/ioutil/tempfile.go?s=1429:1487#L41
	f, err := ioutil.TempFile("", "profile-")
	if err != nil {
		log.Fatalf("Cannot create new temp profile file: %v", err)
	}

	return f
}