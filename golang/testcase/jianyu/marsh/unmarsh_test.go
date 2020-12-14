package marsh

import (
	"encoding/json"
	"testing"
)

func BenchmarkForStruct(b *testing.B) {
	person := InitPerson()
	count := len(person)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ForStruct(person, count)
	}
}

func BenchmarkForRangeStruct(b *testing.B) {
	person := InitPerson()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ForRangeStruct(person)
	}
}

func BenchmarkJsonToStruct(b *testing.B) {
	var (
		person       = InitPerson()
		againPersons []AgainPerson
	)
	data, err := json.Marshal(person)
	if err != nil {
		b.Fatalf("json.Marshal err: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JsonToStruct(data, againPersons)
	}
}
