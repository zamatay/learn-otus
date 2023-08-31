package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"path"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	r := strings.NewReader("test string")
	w := bytes.NewBufferString("")
	t.Run("offset", func(t *testing.T) {
		CopyInternal(r, w, 1, 1025, r.Size())
		require.Equal(t, "est string", w.String())
	})
	t.Run("limit", func(t *testing.T) {
		CopyInternal(r, w, 0, 5, r.Size())
		require.Equal(t, "test ", w.String())
	})
	t.Run("full", func(t *testing.T) {
		CopyInternal(r, w, 0, 0, r.Size())
		require.Equal(t, "test string", w.String())
	})

	fileNmae := path.Join(t.TempDir(), "vks-client02-ios.ovpn.bcp")
	t.Run("work", func(t *testing.T) {
		Copy("~/Загрузки/vks-client02-ios.ovpn", fileNmae, 0, 0)
		require.FileExists(t, fileNmae)
	})
}
