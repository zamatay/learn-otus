package main

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
)

var ErrInvalidOffset = errors.New("InvalidOffset")
var ErrInvalidFileName = errors.New("InvalidFileName")

type File interface {
	fs.File
	io.Seeker
	io.Writer
}

type FileInfo struct {
	file    File
	path    string
	offset  int64
	size    int64
	isWrite bool
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	closeFile := func(file File) {
		err := file.Close()
		if err != nil {
		}
	}
	if fromPath == "" || toPath == "" {
		return ErrInvalidFileName
	}
	fRead := FileInfo{nil, fromPath, offset, 0, false}
	fWrite := FileInfo{nil, toPath, offset, 0, true}
	if isOk, err := getFile(&fRead); !isOk {
		return err
	}
	if isOk, err := getFile(&fWrite); !isOk {
		return err
	}
	defer closeFile(fRead.file)
	defer closeFile(fWrite.file)
	return CopyInternal(fRead.file, fWrite.file, offset, limit, fRead.size)
}

func CopyInternal(src io.Reader, dst io.Writer, o int64, l int64, size int64) error {
	bufSize := getBufferLen(size, o, l)
	_, err := io.CopyN(dst, src, bufSize)
	if err != nil {
		return err
	}
	return nil
}

func getBufferLen(size int64, offset int64, limit int64) int64 {
	if limit > 0 {
		if size-offset > limit {
			return limit
		}
	}
	return size - offset
}

func getFileSize(file fs.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func getFile(fi *FileInfo) (bool, error) {
	var err error
	if fi.isWrite {
		fi.file, err = os.Create(fi.path)
	} else {
		fi.file, err = os.Open(fi.path)
	}
	fi.size, err = getFileSize(fi.file)
	if offset != 0 {
		if fi.size < offset {
			return false, ErrInvalidOffset
		}
		_, err := fi.file.Seek(offset, 0)
		if err != nil {
			return false, err
		}
	}
	if err != nil {
		log.Printf("file name %s error to access", fi.path)
	}
	return true, nil
}
