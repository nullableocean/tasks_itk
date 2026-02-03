package mapstasks

import (
	"maps"
	"testing"
)

func TestWordFrequency(t *testing.T) {
	text := "golang is great and golang is fast"

	expect := map[string]int{
		"golang": 2,
		"is":     2,
		"great":  1,
		"and":    1,
		"fast":   1,
	}

	mfreq := WordFrequency(text)

	if !maps.Equal(expect, mfreq) {
		t.Fatalf("maps nots equal. expect: %v, got: %v", expect, mfreq)
	}
}
func TestWordFrequencyPrint(t *testing.T) {
	text := "golang is great and golang is fast"

	mfreq := WordFrequency(text)
	PrintWordFrequency(mfreq)
}

func BenchmarkPrint(b *testing.B) {
	text := "golang is great and golang is fast"

	mfreq := WordFrequency(text)
	for i := 0; i < b.N; i++ {
		PrintWordFrequency(mfreq)
	}
}
