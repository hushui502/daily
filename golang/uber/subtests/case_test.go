package main

import (
	"github.com/moby/moby/pkg/testutil/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestSubtest(t *testing.T) {
	tests := []struct {
		give     string
		wantHost string
		wantPort string
	}{
		{
			give:     "192.0.2.0:8000",
			wantHost: "192.0.2.0",
			wantPort: "8000",
		},
		{
			give:     "192.0.2.0:http",
			wantHost: "192.0.2.0",
			wantPort: "http",
		},
		{
			give:     ":8000",
			wantHost: "",
			wantPort: "8000",
		},
		{
			give:     "1:8",
			wantHost: "1",
			wantPort: "8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			host, port, err := net.SplitHostPort(tt.give)
			require.NoError(t, err)
			assert.Equal(t, tt.wantHost, host)
			assert.Equal(t, tt.wantPort, port)
		})
	}
}
