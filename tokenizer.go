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

// Splits the input into a paragraph marker (blank line), word, or punctuation.
var rxToken = regexp.MustCompile("(\\s*\\n\\s*\\n\\s*)|([0-9a-zA-Z\\-']+)|([:;,\\.\\?\\!\\-]+)")

func splitTokens(data []byte, atEOF bool) (int, []byte, error) {
	// Handle EOF.
	if atEOF && len(data) == 0 {
		return 0, nil, io.EOF
	}

	// Look for a token.
	if ix := rxToken.FindSubmatchIndex(data); ix != nil {
		// e.g. [22 33 -1 -1 22 33 -1 -1]

		// Paragraph
		if ix[2] >= 0 {
			return ix[3], []byte("\n\n"), nil
		}

		// Word
		if ix[4] >= 0 {
			return ix[5], data[ix[4]:ix[5]], nil
		}

		// Punctuation
		if ix[6] >= 0 {
			return ix[7], data[ix[6]:ix[7]], nil
		}
	}

	// Otherwise, ignore all of this input and move on.
	return len(data), nil, nil
}
