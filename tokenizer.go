package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func Tokenize(in io.Reader) (<-chan string) {
	out := make(chan string, 64)

	scanner := bufio.NewScanner(in)
	scanner.Split(splitTokens)

	go func() {
		for scanner.Scan() {
			out <- scanner.Text()
		}

		if scanner.Err() != nil {
			fmt.Fprintln(os.Stderr, scanner.Err())
			os.Exit(1)
		}

		close(out)
	}()

	return out
}

// Tokenize a word (including punctuation) or a paragraph marker (blank line).
var rxToken = regexp.MustCompile("([#0-9a-zA-Z\\-/!\\.\\?:;,']+)|(\\s*\\n\\s*\\n\\s*)")

func splitTokens(data []byte, atEOF bool) (int, []byte, error) {
	// Handle EOF.
	if atEOF && len(data) == 0 {
		return 0, nil, io.EOF
	}

	// Check for words or paragraph markers, ignore everything else.
	if ix := rxToken.FindSubmatchIndex(data); ix != nil {
		if ix[2] >= 0 {
			// Word
			return ix[3], data[ix[2]:ix[3]], nil
		} else {
			// Paragraph (normalized to a single blank line).
			return ix[5], []byte("\n\n"), nil
		}
	}

	// Otherwise, ignore all of this input and move on.
	return len(data), nil, nil
}
