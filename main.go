package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

var isValidWord = map[string]bool{}
var useOnce = true

func parseWordFile(path string) map[string]bool {
	f, _ := os.Open(path)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		word := strings.ToLower(line)
		isValidWord[word] = true
	}
	return isValidWord
}

func getSubWords(input string) []string {
	subWords := make([]string, 0, 512)
	runes := make([]string, 0, len(input))
	for _, r := range input {
		runes = append(runes, string(r))
	}
	for length := 2; length <= len(input); length += 1 {
		sws := genWordByLength(runes, length)
		subWords = append(subWords, sws...)
	}
	return subWords
}

func genWordByLength(runes []string, l int) []string {
	validWords := make([]string, 0, 10)
	n := numberBaseX{}
	n.initWithLength(l)
	n.base = len(runes)
	chars := make([]string, l)
	for {
		for idx, runeIdx := range n.nums {
			chars[idx] = runes[runeIdx]
		}
		if useOnce && n.hasDuplicate() {
			overflow := n.add1()
			if overflow {
				return validWords
			}
			continue
		}
		word := charsToWord(chars)
		if isValidWord[word] {
			validWords = append(validWords, word)
		}
		overflow := n.add1()
		if overflow {
			return validWords
		}
	}
}

func charsToWord(chars []string) string {
	word := make([]byte, len(chars))
	for idx, r := range chars {
		word[idx] = []byte(r)[0]
	}
	return *(*string)(unsafe.Pointer(&word))
}

// numberBaseX is a fixed length base x number
type numberBaseX struct {
	base int
	nums []int
}

func (n *numberBaseX) initWithLength(l int) {
	n.nums = make([]int, l)
}

func (n *numberBaseX) add1() (overflow bool) {
	for i := len(n.nums) - 1; i >= 0; i-- {
		if n.nums[i] == n.base-1 {
			continue
		}
		n.nums[i] += 1
		if i < len(n.nums)-1 {
			n.nums[i+1] = 0
		}
		return false
	}
	return true
}

// hasDuplicate checks whether there's duplicated num in nums
func (n *numberBaseX) hasDuplicate() bool {
	exist := int64(0)
	for _, num := range n.nums {
		if exist&(1<<num) > 0 {
			return true
		}
		exist |= 1 << num
	}
	return false
}

func main() {
	parseWordFile("./words_alpha.txt")
	results := getSubWords("scotland")
	fmt.Println("when each char could only be used once:\n", results)
	useOnce = false
	results = getSubWords("scotland")
	fmt.Println()
	fmt.Println("when each char could be used multi times:\n", results)
}
