package main

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
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
	fWrite := getFile(toPath, true)
	defer fRead.Close()
	defer fWrite.Close()
	return CopyInternal(fRead, fWrite, offset, limit, getFileSize(fRead))
}

func CopyInternal(read io.ReaderAt, write io.Writer, o int64, l int64, size int64) error {
	if offset >= size {
		return EInvalidOffset
	}
	totalRead := int64(0)
	for {
		bufSize, theEnd := getBufferLen(size, o, totalRead, l)

		data := make([]byte, bufSize)
		cnt, err := read.ReadAt(data, o)
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
	if l > 0 && int64(read)+bufferLen >= l {
		return int(l - int64(read)), true
	}
	if o+bufferLen > size {
		return int(size) - int(o), true
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
		log.Printf("file name %s error to access", path)
	}
	return fRead
}
