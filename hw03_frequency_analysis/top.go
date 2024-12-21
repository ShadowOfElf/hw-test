package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var mask = regexp.MustCompile(`^[^\wа-яА-ЯёЁ]+|[^\wа-яА-ЯёЁ]+$`)

func cleaningWord(s string) string {
	cleanWord := strings.ToLower(s)
	return mask.ReplaceAllString(cleanWord, "")
}

func Top10(inputStr string) []string {
	if len(inputStr) == 0 {
		return nil
	}

	freqMap := map[string]int{}

	for _, word := range strings.Fields(inputStr) {
		cleanWord := cleaningWord(word)
		if cleanWord != "" {
			freqMap[cleanWord]++
		}
	}

	type freqStruct struct {
		word string
		freq int
	}

	resultSlice := make([]freqStruct, len(freqMap))
	i := 0
	for word, count := range freqMap {
		resultSlice[i] = freqStruct{word: word, freq: count}
		i++
	}

	sort.Slice(resultSlice, func(i, j int) bool {
		if resultSlice[i].freq == resultSlice[j].freq {
			return resultSlice[i].word < resultSlice[j].word
		}
		return resultSlice[i].freq > resultSlice[j].freq
	})

	resultStrSlice := make([]string, 0, len(resultSlice))

	for index, freqEl := range resultSlice {
		if index > 9 {
			break
		}
		resultStrSlice = append(resultStrSlice, freqEl.word)
	}

	return resultStrSlice
}
