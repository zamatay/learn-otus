package main

import (
	"github.com/stretchr/testify/require"

	"bytes"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		r := strings.NewReader("test string")
		w := bytes.NewBufferString("")
		err := CopyInternal(r, w, 0, 0, r.Size())
		if err != nil {
			t.Fail()
		}
		require.Equal(t, "test string", w.String())
	})
	t.Run("limit", func(t *testing.T) {
		r := strings.NewReader("test string")
		w := bytes.NewBufferString("")
		err := CopyInternal(r, w, 0, 5, r.Size())
		if err != nil {
			t.Fail()
		}
		require.Equal(t, "test ", w.String())
	})

	//fileTo := path.Join(t.TempDir(), "input.txt")
	//t.Run("real", func(t *testing.T) {
	//	fileFrom := path.Join("./tesdata", "input.txt")
	//	err := Copy(fileFrom, fileTo, 0, 0)
	//	if err != nil {
	//		t.Fail()
	//	}
	//	require.FileExists(t, fileTo)
	//})
}
