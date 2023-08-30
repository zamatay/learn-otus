package main

import (
	"io"
	"io/fs"
	"time"
)

var goodbyeFile = &FakeFile{
	name:     "goodbye",
	contents: "Sayonara, Jim", // 13 bytes, another odd number.
	mode:     0644,
}

// FakeFile implements FileLike and also fs.FileInfo.
type FakeFile struct {
	name     string
	contents string
	mode     fs.FileMode
	offset   int
}

// FileLike methods.

func (f *FakeFile) Name() string {
	// A bit of a cheat: we only have a basename, so that's also ok for FileInfo.
	return f.name
}

func (f *FakeFile) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *FakeFile) Read(p []byte) (int, error) {
	if f.offset >= len(f.contents) {
		return 0, io.EOF
	}
	n := copy(p, f.contents[f.offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) Write(p []byte) (int, error) {
	if f.offset >= len(f.contents) {
		return 0, io.EOF
	}
	Copy(f.contents, p[:])
	n := copy(p, f.contents[f.offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) ReadAt(p []byte, off int64) (int, error) {
	if f.offset >= len(f.contents) {
		return 0, io.EOF
	}
	offset := int64(f.offset) + off
	n := copy(p, f.contents[offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) Close() error {
	return nil
}

// fs.FileInfo methods.

func (f *FakeFile) Size() int64 {
	return int64(len(f.contents))
}

func (f *FakeFile) Mode() fs.FileMode {
	return f.mode
}

func (f *FakeFile) ModTime() time.Time {
	return time.Time{}
}

func (f *FakeFile) IsDir() bool {
	return false
}

func (f *FakeFile) Sys() any {
	return nil
}
