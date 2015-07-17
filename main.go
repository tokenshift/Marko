package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var order int
	switch len(os.Args) {
	case 1:
		order = 1
	case 2:
		i, err := strconv.ParseInt(os.Args[1], 0, 0)

		if err != nil {
			showHelp()
			os.Exit(1)
		}

		if i <= 0 {
			fmt.Fprintln(os.Stderr, "N must be > 0")
			os.Exit(1)
		}

		order = int(i)
	default:
		showHelp()
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Computing Markov chain(s) of order up to", order)

	box := NewMarkoBox()

	//history := make([]string, order, 0)
	previous := ""
	for token := range(inputTokens(os.Stdin)) {
		if previous != "" {
			box.Add(previous, token)
		}
		previous = token
	}

	punctuation := regexp.MustCompile("^\\W+")
	capitalize := true
	for token := range(box.Read()) {
		if !punctuation.MatchString(token) {
			fmt.Print(" ")
		}

		if capitalize {
			fmt.Print(strings.ToUpper(token[0:1]), token[1:])
			capitalize = false
		} else {
			fmt.Print(strings.ToLower(token))
		}

		if token == "." {
			capitalize = true
		}
	}
}

func showHelp() {
	fmt.Fprintln(os.Stderr, "Usage: marko [N] < input.txt")
}

// Breaks an input stream into whitespace-delimited tokens. Any non-word
// characters at the end of a token, which are probably punctuation, is treated
// as its own "word", so that punctuation can also be Markov chained.
func inputTokens(in io.Reader) (<-chan string) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)

	splitWord := regexp.MustCompile("^(\\w+)(\\W*)$")

	out := make(chan string, 1024)
	go func() {
		for scanner.Scan() {
			token := scanner.Text()
			split := splitWord.FindStringSubmatch(token)

			if split == nil {
				continue;
			}

			out <- split[1]
			if split[2] != "" {
				out <- split[2]
			}
		}
		if scanner.Err() != nil {
			fmt.Fprintln(os.Stderr, scanner.Err())
			os.Exit(1)
		}
		close(out)
	}()
	return out
}
