package main

import (
	"math/rand"
)

// Graph of words to successors, with probabilities on each edge.
type MarkoBox struct {
	// The dictionary is used to intern input strings, so that only a single
	// copy of each possible string exists in the graph.
	dictionary map[string]string

	nodes map[string]MarkoNode
}

func NewMarkoBox() MarkoBox {
	return MarkoBox {
		make(map[string]string),
		make(map[string]MarkoNode),
	}
}

type MarkoNode struct {
	edges map[string]int
	total int
}

func (box MarkoBox) intern(input string) string {
	if s, ok := box.dictionary[input]; ok {
		return s
	}

	box.dictionary[input] = input
	return input
}

func (box MarkoBox) Add(precursor, successor string) {
	precursor = box.intern(precursor)
	successor = box.intern(successor)

	if node, ok := box.nodes[precursor]; ok {
		node.Add(successor)
	} else {
		node := MarkoNode { make(map[string]int), 0 }
		node.Add(successor)
		box.nodes[precursor] = node
	}
}

func (node *MarkoNode) Add(successor string) {
	if count, ok := node.edges[successor]; ok {
		node.edges[successor] = count + 1
	} else {
		node.edges[successor] = 1
	}

	node.total += 1
}

// Return a completely random word from the dictionary.
func (box MarkoBox) Random() string {
	r := rand.Intn(len(box.dictionary))
	i := 0
	for s, _ := range(box.dictionary) {
		if i == r {
			return s
		}

		i += 1
	}

	return ""
}

// Return a statistically likely successor to the specified string.
func (box MarkoBox) Next(precursor string) string {
	if node, ok := box.nodes[precursor]; ok {
		r := rand.Intn(node.total)
		for s, count := range(node.edges) {
			r -= count
			if r <= 0 {
				return s
			}
		}
	}

	return box.Random()
}

// A continuous stream of semi-literate gibberish. Attempts to "seed" or
// initialize the stream with a word likely to come after a period.
func (box MarkoBox) Read() (<-chan string) {
	out := make(chan string, 1024)

	go func() {
		token := box.Next(box.intern("."))
		for {
			out <- token
			token = box.Next(token)
		}
	}()

	return out
}
