package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

var (
	nf Interface

	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	info    = color.New(color.FgGreen).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgHiMagenta).SprintFunc()
)

func start() {
	var err error

	nf, err = New()
	if err != nil {
		log.Fatal(err)
	}

	var (
		recentRankLimit = 10

		sigch   = make(chan os.Signal, 1)
		ticker  = time.NewTicker(3 * time.Second)
		timeout = time.NewTimer(300 * time.Second)
	)

	signal.Notify(sigch, syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)

	defer func() {
		nf.Stop()
	}()

	for {
		select {
		case <-sigch:
			return

		case <-timeout.C:
			return

		case <-ticker.C:
			rank, err := nf.GetProcessRank(recentRankLimit, 3)
			if err != nil {
				log.Errorf("GetProcessRank failed, err: %s", err.Error())
				continue
			}

			clear()
			showTable(rank)
			time.Sleep(5 * time.Second)
		}
	}
}

func stop() {
	if nf == nil {
		return
	}
	nf.Stop()
}

const thold = 1024 * 1024 // 1mb

func clear() {
	fmt.Printf("\x1b[2J")
}

func showTable(ps []*Process) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"pid", "name", "exe", "inodes", "sum_in", "sum_out", "in_rate", "out_rate"})
	table.SetRowLine(true)

	items := [][]string{}
	for _, po := range ps {
		inRate := humanBytes(po.TrafficStats.InRate)
		if po.TrafficStats.InRate > int64(thold) {
			inRate = red(inRate)
		}

		outRate := humanBytes(po.TrafficStats.OutRate)
		if po.TrafficStats.OutRate > int64(thold) {
			outRate = red(outRate)
		}

		item := []string{
			po.Pid,
			po.Name,
			po.Exe,
			cast.ToString(po.InodeCount),
			humanBytes(po.TrafficStats.In),
			humanBytes(po.TrafficStats.Out),
			inRate,
			outRate,
		}

		items = append(items, item)
	}

	table.AppendBulk(items)
	table.Render()
}

func humanBytes(n int64) string {
	return humanize.Bytes(uint64(n))
}

func main() {
	log.Info("start netflow sniffer")

	start()
	stop()

	log.Info("netflow sniffer exit")
}
