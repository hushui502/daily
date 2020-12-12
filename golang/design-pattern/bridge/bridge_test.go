package bridge

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrNotification_Notify(t *testing.T) {
	sender := NewEmailMsgSender([]string{"hff@gmail.com"})
	n := NewErrNotification(sender)
	err := n.Notify("server cash down")

	assert.Nil(t, err)
}