package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const (
	bufferLen = 1024
)

var EInvalidOffset = errors.New("invalid offset")

type Stat interface {
	Stat() (fs.FileInfo, error)
}

type ReadAtWriteCloser interface {
	io.WriteCloser
	io.ReaderAt
	Stat
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fRead := getFile(fromPath, false)
	defer fRead.Close()
	fWrite := getFile(toPath, true)
	defer fWrite.Close()
	return CopyInternal(fRead, fWrite, offset, limit)
}

func CopyInternal(read ReadAtWriteCloser, write ReadAtWriteCloser, o int64, l int64) error {

	size := getFileSize(read)

	if offset >= size {
		return EInvalidOffset
	}
	totalRead := int64(0)
	for {
		bufSize, theEnd := getBufferLen(size, offset, totalRead, limit)

		data := make([]byte, bufSize)
		cnt, err := read.ReadAt(data, offset)
		if err != nil {
			return err
		}
		offset += int64(cnt)
		totalRead += int64(cnt)
		write.Write(data)
		if theEnd {
			break
		}
	}

	return nil
}

func getBufferLen(size int64, o int64, read int64, l int64) (int, bool) {
	if o+bufferLen > size {
		return int(size) - int(o), true
	}
	if l > 0 && int64(read)+bufferLen >= l {
		return int(l - int64(read)), true
	}
	return bufferLen, false
}

func getFileSize(file ReadAtWriteCloser) int64 {
	var size64 int64
	if info, err := file.Stat(); err == nil {
		size64 = info.Size()
	}
	return size64
}

func getFile(path string, isWrite bool) ReadAtWriteCloser {
	var fRead *os.File
	var err error
	if isWrite {
		fRead, err = os.Create(path)
	} else {
		fRead, err = os.Open(path)
	}
	if err != nil {
		fmt.Errorf("file name %s error to access", path)
	}
	return fRead
}
