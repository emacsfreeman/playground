package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	var order = 3
	var hash = "#"
	carpet := []string{hash}
	for ; order > 0; order-- {
		hole := strings.Repeat(" ", utf8.RuneCountInString(carpet[0]))
		middle := make([]string, len(carpet))
		for i, s := range carpet {
			middle[i] = s + hole + s
			carpet[i] = strings.Repeat(s, 3)
		}
		carpet = append(append(carpet, middle...), carpet...)
	}

	for _, r := range carpet {
		fmt.Println(r)
	}
}
