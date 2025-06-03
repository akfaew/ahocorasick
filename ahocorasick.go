package goahocorasick

import (
	"os"

	godarts "github.com/akfaew/darts"

	"fmt"
)

const FAIL_STATE = -1
const ROOT_STATE = 1

type Machine struct {
	trie    *godarts.DoubleArrayTrie
	failure []int
	output  map[int]([][]byte)
}

type Term struct {
	Pos  int
	Word []byte
}

func (m *Machine) Build(keywords [][]byte) (err error) {
	if len(keywords) == 0 {
		return fmt.Errorf("empty keywords")
	}

	d := new(godarts.Darts)

	var trie *godarts.LinkedListTrie
	m.trie, trie, err = d.Build(keywords)
	if err != nil {
		return err
	}

	m.output = make(map[int]([][]byte))
	for idx, val := range d.Output {
		m.output[idx] = append(m.output[idx], val)
	}

	queue := make([](*godarts.LinkedListTrieNode), 0)
	m.failure = make([]int, len(m.trie.Base))
	for _, c := range trie.Root.Children {
		if c.Base == -1 {
			for _, v := range keywords {
				fmt.Fprintf(os.Stderr, "keyword: %s\n", string(v))
			}
			return fmt.Errorf("invalid trie (c.Base == -1). len(keywords) = %d", len(keywords)) // to debug why sometimes it's -1
		}
		m.failure[c.Base] = godarts.ROOT_NODE_BASE
	}
	queue = append(queue, trie.Root.Children...)

	for {
		if len(queue) == 0 {
			break
		}

		node := queue[0]
		for _, n := range node.Children {
			if n.Base == godarts.END_NODE_BASE {
				continue
			}
			inState := m.f(node.Base)
		set_state:
			outState := m.g(inState, n.Code-godarts.ROOT_NODE_BASE)
			if outState == FAIL_STATE {
				inState = m.f(inState)
				goto set_state
			}
			if _, ok := m.output[outState]; ok {
				copyOutState := make([][]byte, 0)
				copyOutState = append(copyOutState, m.output[outState]...)
				m.output[n.Base] = append(copyOutState, m.output[n.Base]...)
			}
			m.setF(n.Base, outState)
		}
		queue = append(queue, node.Children...)
		queue = queue[1:]
	}

	return nil
}

func (m *Machine) g(inState int, input byte) (outState int) {
	if inState == FAIL_STATE {
		return ROOT_STATE
	}

	t := inState + int(input) + godarts.ROOT_NODE_BASE
	if t >= len(m.trie.Base) {
		if inState == ROOT_STATE {
			return ROOT_STATE
		}
		return FAIL_STATE
	}
	if inState == m.trie.Check[t] {
		return m.trie.Base[t]
	}

	if inState == ROOT_STATE {
		return ROOT_STATE
	}

	return FAIL_STATE
}

func (m *Machine) f(index int) (state int) {
	return m.failure[index]
}

func (m *Machine) setF(inState, outState int) {
	m.failure[inState] = outState
}

func (m *Machine) MultiPatternSearch(content []byte, returnImmediately bool) [](*Term) {
	terms := make([](*Term), 0)

	state := ROOT_STATE
	for pos, c := range content {
	start:
		if m.g(state, c) == FAIL_STATE {
			state = m.f(state)
			goto start
		} else {
			state = m.g(state, c)
			if val, ok := m.output[state]; ok {
				for _, word := range val {
					term := new(Term)
					term.Pos = pos - len(word) + 1
					term.Word = word
					terms = append(terms, term)
					if returnImmediately {
						return terms
					}
				}
			}
		}
	}

	return terms
}

func (m *Machine) MultiPatternSearchQuick(content []byte) (ret []string) {
	state := ROOT_STATE
	tmp := map[string]bool{}
	for _, c := range content {
	start:
		if m.g(state, c) == FAIL_STATE {
			state = m.f(state)
			goto start
		} else {
			state = m.g(state, c)
			if val, ok := m.output[state]; ok {
				for _, word := range val {
					tmp[string(word)] = true
				}
			}
		}
	}

	for k := range tmp {
		ret = append(ret, k)
	}
	return ret
}

func (m *Machine) ExactSearch(content []byte) [](*Term) {
	if m.trie.ExactMatchSearch(content, 0) {
		t := new(Term)
		t.Word = content
		t.Pos = 0
		return [](*Term){t}
	}

	return nil
}
