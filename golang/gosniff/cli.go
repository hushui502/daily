package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.6.0"

func NewApp() *cobra.Command {
	defaultOpts := DefaultOptions()

	opt := Options{}
	var (
		mode int
		unit string
		list bool
	)

	app := &cobra.Command{
		Use:     "sniffer",
		Short:   "# A modern alternative network traffic sniffer.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			if list {
				devices, err := ListAllDevices()
				if err != nil {
					exit(err.Error())
				}
				for _, device := range devices {
					fmt.Println(device.Name)
				}
				return
			}
			opt.ViewMode = ViewMode(mode)
			opt.Unit = Unit(unit)
			if err := opt.Validate(); err != nil {
				exit(err.Error())
			}

			sniffer, err := NewSniffer(opt)
			if err != nil {
				exit(err.Error())
			}
			defer sniffer.Close()
			sniffer.Start()
		},
		Example: `  # processes mode for pid 1024,2048 in MB unit
  $ sniffer -p 1024 -p 2048 -m 2 -u MB
  # only capture the TCP protocol packets with lo,eth prefixed devices
  $ sniffer -b tcp -d lo -d eth`,
	}

	app.Flags().BoolVarP(&list, "list", "l", false, "list all devices name")
	app.Flags().BoolVarP(&opt.AllDevices, "all-devices", "a", false, "listen all devices if present")
	app.Flags().StringVarP(&opt.BPFFilter, "bpf", "b", defaultOpts.BPFFilter, "specify string pcap filter with the BPF syntax")
	app.Flags().IntVarP(&opt.Interval, "interval", "i", defaultOpts.Interval, "interval for refresh rate in seconds")
	app.Flags().StringArrayVarP(&opt.DevicesPrefix, "devices-prefix", "d", defaultOpts.DevicesPrefix, "prefixed devices to monitor")
	app.Flags().BoolVarP(&opt.DisableDNSResolve, "no-dns-resolve", "n", defaultOpts.DisableDNSResolve, "disable the DNS resolution")
	app.Flags().Int32SliceVarP(&opt.Pids, "pids", "p", defaultOpts.Pids, "pids to watch, empty stands for all pids")
	app.Flags().IntVarP(&mode, "mode", "m", int(defaultOpts.ViewMode), "view mode of sniffer (0: bytes 1: packets 2: processes)")
	app.Flags().StringVarP(&unit, "unit", "u", defaultOpts.Unit.String(), "unit of traffic stats, optional: B, Kb, KB, Mb, MB, Gb, GB")

	app.Flags().PrintDefaults()
	return app
}

func main() {
	app := NewApp()
	if err := app.Execute(); err != nil {
		exit(err.Error())
	}
}
