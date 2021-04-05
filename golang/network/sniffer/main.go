package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

var (
	ctx, cancelFunc = context.WithCancel(context.Background())
	result          = make(chan string, 1)
)

func main() {
	iface := flag.String("i", "", "interface capture")
	flag.Parse()
	flag.Usage = usage

	if len(os.Args) < 2 || len(os.Args[2]) == 0 {
		usage()
		os.Exit(-1)
	}

	go func() {
		for res := range result {
			fmt.Fprintln(os.Stdout, res)
		}
	}()

	sniffer := NewSniffer(
		SetDevice(*iface),
		SetCtx(ctx),
		SetSyncCh(result))

	if err := sniffer.Start(); err != nil {
		log.Fatal(err)
	}

	defer sniffer.Stop()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	cancelFunc()

	fmt.Fprintln(os.Stdout, "End sniff :)")
}
