package main

import (
	"fmt"
	"strings"
)

func wordCount(str string) map[string]int {
	wordList := strings.Fields(str)
	counts := make(map[string]int)
	for _, word := range wordList {
		_, ok := counts[word]
		if ok {
			counts[word] += 1
		} else {
			counts[word] = 1
		}
	}
	return counts
}

func main() {
	strLine := "Australia Canada Germany Australia Japan Canada"
	for index, element := range wordCount(strLine) {
		fmt.Println(index, "=>", element)
	}
}
