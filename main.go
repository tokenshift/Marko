package main

import (
	"fmt"
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

	previous := ""
	for token := range(Tokenize(os.Stdin)) {
		if previous != "" {
			box.Add(previous, token)
		}
		previous = token
	}

	capitalize := true
	noSpace := true

	rxPunctuation    := regexp.MustCompile("^\\W+")
	rxCapitalizeNext := regexp.MustCompile("([\\.!\\?]|\\n\\s*\\n\\s*)$")
	rxNoSpaceNext    := regexp.MustCompile("(\\-|\\n\\s*\\n\\s*)$")

	for word := range box.Read() {
		if !noSpace && !rxPunctuation.MatchString(word) {
			fmt.Print(" ")
		}

		if capitalize {
			fmt.Print(strings.ToUpper(word[0:1]), word[1:])
		} else {
			fmt.Print(word)
		}

		capitalize = rxCapitalizeNext.MatchString(word)
		noSpace = rxNoSpaceNext.MatchString(word)
	}
}

func showHelp() {
	fmt.Fprintln(os.Stderr, "Usage: marko [N] < input.txt")
}
