package goahocorasick

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/akfaew/utils/fixture"
	"github.com/stretchr/testify/require"
)

func Read(t *testing.T, filename string) ([][]byte, error) {
	t.Helper()

	dict := [][]byte{}

	f, err := os.OpenFile("testdata/input/"+filename, os.O_RDONLY, 0660)
	require.NoError(t, err)

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadBytes('\n')
		if err != nil || err == io.EOF {
			break
		}
		l = bytes.TrimSpace(l)
		dict = append(dict, l)
	}

	return dict, nil
}

func TestBuild(t *testing.T) {
	keywords, err := Read(t, "keywords")
	require.NoError(t, err)

	m := new(Machine)
	require.NoError(t, m.Build(keywords))
}

func TestMultiPatternSearch(t *testing.T) {
	keywords, err := Read(t, "keywords")
	require.NoError(t, err)
	m := new(Machine)
	require.NoError(t, m.Build(keywords))

	content := []byte("ushers she she")
	terms := m.MultiPatternSearch(content, false)
	fixture.Fixture(t, terms)
}

func TestMultiPatternSearchQuick(t *testing.T) {
	keywords, err := Read(t, "keywords")
	require.NoError(t, err)
	m := new(Machine)
	require.NoError(t, m.Build(keywords))

	content := []byte("ushers she she")
	ret := m.MultiPatternSearchQuick(content)
	fixture.Fixture(t, ret)
}

func TestExactSearch(t *testing.T) {
	keywords, err := Read(t, "keywords")
	require.NoError(t, err)
	m := new(Machine)
	require.NoError(t, m.Build(keywords))

	for _, k := range keywords {
		if m.ExactSearch(k) == nil {
			t.Error("exact search failed")
		}
	}
	require.Len(t, keywords, 4)
}
