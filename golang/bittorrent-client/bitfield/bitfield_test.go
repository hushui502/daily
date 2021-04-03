package bitfield

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitfield_HasPiece(t *testing.T) {
	bf := Bitfield{0b01010100, 0b01010100}
	outputs := []bool{false, true, false, true, false, true, false, false, false, true, false, true, false, true, false, false, false, false, false, false}
	for i := 0; i < len(outputs); i++ {
		assert.Equal(t, outputs[i], bf.HasPiece(i))
	}
}

func TestBitfield_SetPiece(t *testing.T) {
	tests := []struct {
		input  Bitfield
		index  int
		output Bitfield
	}{
		{
			input:  Bitfield{0b01010100, 0b01010100},
			index:  4, //          v (set)
			output: Bitfield{0b01011100, 0b01010100},
		},
		{
			input:  Bitfield{0b01010100, 0b01010100},
			index:  9, //                   v (noop)
			output: Bitfield{0b01010100, 0b01010100},
		},
		{
			input:  Bitfield{0b01010100, 0b01010100},
			index:  15, //                        v (set)
			output: Bitfield{0b01010100, 0b01010101},
		},
		{
			input:  Bitfield{0b01010100, 0b01010100},
			index:  19, //                            v (noop)
			output: Bitfield{0b01010100, 0b01010100},
		},
	}

	for _, item := range tests {
		bf := item.input
		bf.SetPiece(item.index)
		assert.Equal(t, item.output, bf)
	}
}
