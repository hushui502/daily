package sensor

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go-sec/analysis/sensor/misc"
	"go-sec/analysis/sensor/settings"
	"time"
)

var (
	device      string
	snapshotlen int32 = 1024
	promiscuous bool  = true
	err         error
	timeout     time.Duration = pcap.BlockForever
	handle      *pcap.Handle

	DebugMode bool
	filter    = ""

	ApiUrl    string
	SecureKey string
)

func init() {
	device = settings.DeviceName
	DebugMode = settings.DebugMode
	filter = settings.FilterRule

	cfg := settings.Cfg
	serverSec := cfg.Section("server")
	ApiUrl = serverSec.Key("API_URL").MustString("")
	SecureKey = serverSec.Key("API_KEY").MustString("")
}

func Start(ctx *cli.Context) {
	if ctx.IsSet("debug") {
		DebugMode = ctx.Bool("debug")
	}
	if DebugMode {
		misc.Log.Logger.Level = logrus.DebugLevel
	}

	if ctx.IsSet("length") {
		snapshotlen = int32(ctx.Int("len"))
	}

	handle, err := pcap.OpenLive(device, snapshotlen, promiscuous, timeout)
	if err != nil {
		misc.Log.Fatal(err)
	}
	defer handle.Close()

	if ctx.IsSet("filter") {
		filter = ctx.String("filter")
	}
	handle.SetBPFFilter(filter)
	misc.Log.Infof("set SetBPFFilter: %v, err: %v", filter, err)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	ProcessPackets(packetSource.Packets())
}
