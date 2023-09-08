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

var ErrInvalidOffset = errors.New("invalid offset")

type Stat interface {
	Stat() (fs.FileInfo, error)
}

type ReadAtWriteCloser interface {
	io.WriteCloser
	io.ReaderAt
	Stat
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fRead, size := getFile(fromPath, false)
	fWrite, _ := getFile(toPath, true)
	defer fRead.Close()
	defer fWrite.Close()
	return CopyInternal(fRead, fWrite, offset, limit, size)
}

func CopyInternal(src io.ReaderAt, dst io.Writer, o int64, l int64, size int64) error {
	if o >= size {
		return ErrInvalidOffset
	}
	totalRead := int64(0)
	for {
		bufSize, theEnd := getBufferLen(size, o, totalRead, l)

		//io.CopyN(dst, src, bufSize)
		data := make([]byte, bufSize)
		cnt, err := src.ReadAt(data, o)
		if err != nil {
			return err
		}
		o += int64(cnt)
		totalRead += int64(cnt)
		dst.Write(data)
		if theEnd {
			break
		}
	}
	return nil
}

func getBufferLen(size int64, o int64, read int64, l int64) (int64, bool) {
	if l > 0 && int64(read)+bufferLen >= l {
		return l - read, true
	}
	if o+bufferLen > size {
		return size - o, true
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

func getFile(path string, isWrite bool) (ReadAtWriteCloser, int64) {
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
	return fRead, getFileSize(fRead)
}
