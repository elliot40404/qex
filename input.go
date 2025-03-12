package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func input(text string) []any {
	re := regexp.MustCompile(`\$(\d+)`)
	matches := re.FindAllStringSubmatch(text, -1)
	uniqueKeys := make(map[string]bool)
	for _, match := range matches {
		uniqueKeys[match[1]] = true
	}
	var keys []int
	for k := range uniqueKeys {
		num, _ := strconv.Atoi(k)
		keys = append(keys, num)
	}
	userInputs := make(map[int]string)
	scanner := bufio.NewScanner(os.Stdin)
	for _, k := range keys {
		fmt.Printf("Enter value for $%d: ", k)
		scanner.Scan()
		userInputs[k] = scanner.Text()
	}
	results := make([]any, len(keys))
	for i, k := range keys {
		results[i] = userInputs[k]
	}
	return results
}
