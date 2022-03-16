//go:build windows
// +build windows

package hostfile

import (
	"os"
)

var HostsPath = os.Getenv("SystemRoot") + `\System32\drivers\etc\hosts`
