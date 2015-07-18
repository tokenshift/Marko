package main

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
)

// Graph of words to successors, with probabilities on each edge.
type MarkoBox struct {
	order int

	// The dictionary is used to intern input strings, so that only a single
	// copy of each possible string exists in the graph.
	dictionary map[string]string

	nodes map[string]MarkoNode
}

func NewMarkoBox(order int) MarkoBox {
	return MarkoBox {
		order,
		make(map[string]string),
		make(map[string]MarkoNode),
	}
}

type MarkoNode struct {
	edges map[string]int
}

func (box *MarkoBox) intern(input string) string {
	if s, ok := box.dictionary[input]; ok {
		return s
	}

	box.dictionary[input] = input
	return input
}

func (box *MarkoBox) Add(precursors []string, successor string) {
	if len(precursors) > box.order {
		precursors = precursors[len(precursors) - box.order:]
	}

	for i := 0; i < len(precursors); i += 1 {
		precursor := strings.Join(precursors[i:], " ")
		precursor = box.intern(precursor)
		successor = box.intern(successor)

		if node, ok := box.nodes[precursor]; ok {
			node.Add(successor)
		} else {
			node := MarkoNode { make(map[string]int) }
			node.Add(successor)
			box.nodes[precursor] = node
		}
	}
}

func (node *MarkoNode) Add(successor string) {
	if count, ok := node.edges[successor]; ok {
		node.edges[successor] = count + 1
	} else {
		node.edges[successor] = 1
	}
}

// Return a completely random word from the dictionary.
func (box *MarkoBox) Random() string {
	fmt.Print("HERE")
	r := rand.Intn(len(box.dictionary))
	i := 0
	for s, _ := range(box.dictionary) {
		if i == r {
			return s
		}

		i += 1
	}

	panic("Something is broken.")
}

// Return a statistically likely successor to the specified string.
func (box *MarkoBox) Next(precursors []string) string {
	for i := 0; i < len(precursors); i += 1 {
		precursor := strings.Join(precursors[i:], " ")

		if node, ok := box.nodes[precursor]; ok {
			r := rand.Intn(node.Total())
			for s, count := range(node.edges) {
				r -= count
				if r <= 0 {
					return s
				}
			}
		}
	}

	return box.Random()
}

// A continuous stream of semi-literate gibberish. Attempts to "seed" or
// initialize the stream with a word likely to come after a period.
func (box *MarkoBox) Read() (<-chan string) {
	out := make(chan string, 1024)

	go func() {
		precursors := []string{"\n\n"}

		token := box.Next(precursors)
		for {
			out <- token

			if len(precursors) >= box.order {
				precursors = append(precursors[1:], token)
			} else {
				precursors = append(precursors, token)
			}

			token = box.Next(precursors)
		}
	}()

	return out
}

func (node MarkoNode) Total() int {
	sum := 0
	for _, count := range node.edges {
		sum += count
	}
	return sum
}

func (box MarkoBox) Print(out io.Writer) {
	for precursor, node := range box.nodes {
		if precursor == "\n\n" {
			precursor = "{P}"
		}

		for successor, count := range node.edges {
			if successor == "\n\n" {
				successor = "{P}"
			}

			fmt.Fprintf(out, "%s => %s [%d/%d]\n", precursor, successor, count, node.Total())
		}
		fmt.Fprintln(out)
	}
}
