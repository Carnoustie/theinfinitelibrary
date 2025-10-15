package main

import (
	"fmt"
	"strings"
)

func main() {
	st := "goodevening"
	fmt.Println(st[0 : len(st)-1])

	s1 := "Anna Karenina"
	s2 := "War and Peace"
	s1 = strings.ReplaceAll(s1, " ", "a")
	s2 = strings.ReplaceAll(s2, " ", "a")
	l1 := len(s1)
	l2 := len(s2)
	res := levenshteinDistance(s1, s2, l1, l2)
	fmt.Println("\n\nres: ", res)
}

// func levenshteinDistance(s1 string, s2 string, len_s1 int, len_s2 int) int {
// 	if len_s1 == 0 {
// 		return len_s2
// 	}
// 	return len_s1
// 	// else {
// 	// 	return 1 + levenshteinDistance(s1, s2, len_s1-1, len_s2-1)
// 	// }
// }

func levenshteinDistance(s1 string, s2 string, len_s1 int, len_s2 int) int {
	if len_s1 == 0 {
		return len_s2
	} else if len_s2 == 0 {
		return len_s1
	} else if len_s1 == len_s2 && s1[len_s1-1] == s2[len_s2-1] {
		return levenshteinDistance(s1, s2, len_s1-1, len_s2-1)
	} else {
		return 1 + min(levenshteinDistance(s1, s2, len_s1, len_s2-1), min(levenshteinDistance(s1, s2, len_s1-1, len_s2), levenshteinDistance(s1, s2, len_s1-1, len_s2-1)))
	}
}
