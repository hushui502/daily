package golib

import (
	"sync"
)

var loadSync sync.Once

func main() {
}

func loadImage() {
}

func load() {
	loadSync.Do(loadImage)
}
