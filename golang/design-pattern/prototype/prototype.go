package prototype

import (
	"encoding/json"
	"time"
)

type Keyword struct {
	word string
	visit int
	UpdatedAt *time.Time
}

func (k *Keyword) Clone() *Keyword {
	var newKeyWord Keyword
	b, _ := json.Marshal(k)
	json.Unmarshal(b, &newKeyWord)

	return &newKeyWord
}

type Keywords map[string]*Keyword

func (words Keywords) Clone(updateWords []*Keyword) Keywords {
	newKeywords := Keywords{}

	for k, v := range words {
		newKeywords[k] = v
	}

	for _, word := range updateWords {
		newKeywords[word.word] = word.Clone()
	}

	return newKeywords
}
