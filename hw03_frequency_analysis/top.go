package hw03frequencyanalysis

import (
	"container/heap"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var taskWithAsteriskIsCompleted = true

type wordCount struct {
	word  string
	count int
}

type wordHeap []wordCount

func (h wordHeap) Len() int {
	return len(h)
}

func (h wordHeap) Less(i, j int) bool {
	if h[i].count == h[j].count {
		return h[i].word > h[j].word
	}
	return h[i].count < h[j].count
}

func (h wordHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *wordHeap) Push(x interface{}) {
	*h = append(*h, x.(wordCount))
}

func (h *wordHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func MeasureTime(fn func(text string) []string, text string) []string {
	start := time.Now()

	res := fn(text)

	elapsed := time.Since(start)
	fmt.Println("Время выполнения функции:", elapsed)
	return res
}

// Top10 returns the 10 most frequently occurring words.
func Top10(text string) []string {
	parseWords := splitWithOrWithoutAsterisk(taskWithAsteriskIsCompleted, text)
	counters := make(map[string]int)

	for _, word := range parseWords {
		if taskWithAsteriskIsCompleted && word == "-" {
			continue
		}
		counters[word]++
	}

	wh := &wordHeap{}
	heap.Init(wh)
	for word, count := range counters {
		heap.Push(wh, wordCount{word, count})
		if wh.Len() > 10 {
			heap.Pop(wh)
		}
	}

	result := make([]string, 0, 10)
	for wh.Len() > 0 {
		word := heap.Pop(wh).(wordCount).word
		result = append([]string{word}, result...)
	}

	return result
}

// Top10Old returns the 10 most frequently occurring words.
func Top10Old(text string) []string {
	parseWords := splitWithOrWithoutAsterisk(taskWithAsteriskIsCompleted, text)
	counters := make(map[string]int)

	for _, word := range parseWords {
		if taskWithAsteriskIsCompleted && word == "-" {
			continue
		}
		counters[word]++
	}

	words := make([]wordCount, 0, len(counters))
	for word, count := range counters {
		words = append(words, wordCount{word, count})
	}

	sort.Slice(words, func(i, j int) bool {
		if words[i].count == words[j].count {
			return words[i].word < words[j].word
		}
		return words[i].count > words[j].count
	})

	result := make([]string, 0, 10)
	for i, wordCount := range words {
		if i >= 10 {
			break
		}
		result = append(result, wordCount.word)
	}

	return result
}

// p{L} - any letters of any language.
// p{N} - any kind of numeric characters in any language.
var wordRegex = regexp.MustCompile(`[\p{L}\p{N}-]+`)

func splitWithOrWithoutAsterisk(sign bool, text string) []string {
	if !sign {
		return strings.Fields(text)
	}

	text = strings.ToLower(text)
	words := wordRegex.FindAllString(text, -1)
	return words
}
