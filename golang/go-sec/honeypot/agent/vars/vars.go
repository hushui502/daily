package vars

import "github.com/orcaman/concurrent-map"

var (
	ProxyHost string
	ProxyPort int
	DebugMode bool

	CurrentDir string
	CaCert     string
	CaKey      string

	Cmap = cmap.New()
)
