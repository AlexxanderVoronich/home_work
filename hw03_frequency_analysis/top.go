package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

// Top10 returns the 10 most frequently occurring words.
func Top10(text string) []string {
	parseWords := splitWithOrWithoutAsterisk(false, text)
	counters := make(map[string]int)

	for _, word := range parseWords {
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

	var result []string
	for _, word := range words {
		if word != "-" {
			result = append(result, word)
		}
	}

	return result
}
