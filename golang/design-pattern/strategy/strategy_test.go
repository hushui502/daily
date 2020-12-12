package strategy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileStorage_Save(t *testing.T) {
	data := []byte("hello")
	strategyType := "file"
	//if sensitive {
	//	strategyType = "encrypt_file"
	//}

	storage, err := NewStorageStrategy(strategyType)
	assert.NoError(t, err)
	assert.NoError(t, storage.Save("./test.txt", data))
}
