package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

var rxParagraph = regexp.MustCompile("^\\s*\\n\\s*\\n\\s*$")

func Tokenize(in io.Reader) (<-chan string) {
	out := make(chan string, 64)

	scanner := bufio.NewScanner(in)
	scanner.Split(splitTokens)

	go func() {
		for scanner.Scan() {
			// Normalize paragraphs.
			if rxParagraph.MatchString(scanner.Text()) {
				out <- "\n\n"
			} else {
				out <- scanner.Text()
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

// Splits the input into a paragraph marker (blank line), word, or punctuation.
var rxToken = regexp.MustCompile("(\\s*\\n\\s*\\n\\s*)|([0-9a-zA-Z\\-%'\\.]+)|([\\.!\\?\\-]+)")

func splitTokens(data []byte, atEOF bool) (int, []byte, error) {
	// Handle EOF.
	if atEOF && len(data) == 0 {
		return 0, nil, io.EOF
	}

	// Look for a token.
	if ix := rxToken.FindIndex(data); ix != nil {
		return ix[1], data[ix[0]:ix[1]], nil
	}

	// Otherwise, ignore all of this input and move on.
	return len(data), nil, nil
}
