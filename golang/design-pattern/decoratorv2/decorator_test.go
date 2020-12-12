package decoratorv2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColorSquare_Draw(t *testing.T) {
	sq := Square{}
	csq := NewColorSquare(sq, "blue")
	got := csq.Draw()

	assert.Equal(t, "this is a square, color is blue", got)
}
