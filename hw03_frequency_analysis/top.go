package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var regex = regexp.MustCompile("[\\p{Cyrillic}\\w]+(-[\\p{Cyrillic}\\w]+)*")

type WordCount struct {
	Word  string
	Count int
}

func Top10(text string) []string {
	words := strings.Fields(text)
	wordsCountMap := make(map[string]WordCount)

	for _, word := range words {
		if taskWithAsteriskIsCompleted {
			word = strings.ToLower(regex.FindString(word))
			if word == "" {
				continue
			}
		}

		if wordCount, ok := wordsCountMap[word]; ok {
			wordCount.Count++
			wordsCountMap[word] = wordCount
		} else {
			wordsCountMap[word] = WordCount{word, 1}
		}
	}

	wordsCountSlice := make([]WordCount, 0, len(wordsCountMap))

	for _, value := range wordsCountMap {
		wordsCountSlice = append(wordsCountSlice, value)
	}

	sort.Slice(wordsCountSlice, func(i, j int) bool {
		if wordsCountSlice[i].Count != wordsCountSlice[j].Count {
			return wordsCountSlice[i].Count > wordsCountSlice[j].Count
		}

		return wordsCountSlice[i].Word < wordsCountSlice[j].Word
	})

	top10Words := make([]string, 0, 10)
	for i, value := range wordsCountSlice {
		if i >= 10 {
			break
		}

		top10Words = append(top10Words, value.Word)
	}

	return top10Words
}
