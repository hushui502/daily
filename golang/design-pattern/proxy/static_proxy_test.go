package proxy

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewUserProxy(t *testing.T) {
	proxy := NewUserProxy(&User{})

	err := proxy.Login("test", "testpassword")

	require.Nil(t, err)
}
