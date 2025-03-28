package main

import (
	"fmt"
	"strings"
)

func WordCount(s string) map[string]int {
	ans := make(map[string]int)
	for _, word := range strings.Fields(s) {
		ans[word]++
	}
	return ans
}

func main() {
	//wc.Test(WordCount)
	fmt.Println(WordCount("I am learning GO! I am writing code in GO!"))

}
