package main

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
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
			log.Printf("%v", err)
		}
	}
	if fromPath == "" || toPath == "" {
		return ErrInvalidFileName
	}
	fRead := FileInfo{nil, fromPath, offset, 0, false}
	fWrite := FileInfo{nil, toPath, 0, 0, true}
	if isOk, err := getFile(&fRead); !isOk {
		log.Printf("%v", err)
		return err
	}
	if isOk, err := getFile(&fWrite); !isOk {
		log.Printf("%v", err)
		return err
	}
	defer closeFile(fRead.file)
	defer closeFile(fWrite.file)
	return CopyInternal(fRead.file, fWrite.file, offset, limit, fRead.size)
}

func CopyInternal(src io.Reader, dst io.Writer, o int64, l int64, size int64) error {
	if o < 0 {
		return ErrInvalidOffset
	}
	bufSize := getBufferLen(size, o, l)
	bar := pb.Full.Start64(bufSize)
	barReader := bar.NewProxyReader(src)
	_, err := io.CopyN(dst, barReader, bufSize)
	bar.Finish()
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func getBufferLen(size int64, offset int64, limit int64) int64 {
	if limit > 0 && size-offset > limit {
		return limit
	}
	return size - offset
}

func getFileSize(file fs.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		log.Printf("%v", err)
		return 0, err
	}
	return info.Size(), nil
}

func getFile(fi *FileInfo) (bool, error) {
	var err error
	if fi.isWrite {
		fi.file, err = os.Create(fi.path)
		if err != nil {
			log.Printf("%v", err)
			return false, err
		}
		return true, nil
	}

	fi.file, err = os.Open(fi.path)
	if err != nil {
		log.Printf("%v", err)
		return false, err
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
