package main

import (
	"testing"
)

//type StorageMock struct {
//}
//
//func (s StorageMock) Close() error {
//	return nil
//}
//
//func (s StorageMock) ReadAt(p []byte, off int64) (n int, err error){
//	return 0, nil
//}
//
//func (s StorageMock) Write(p []byte) (n int, err error){
//	return 0, nil
//}
//
//func (s StorageMock) Stat() (fs.FileInfo, error){
//	return {}
//}

var testFile = &FakeFile{
	name:     "goodbye",
	contents: "Это тестовая строка для копирования", // 13 bytes, another odd number.
	mode:     0644,
}

func TestCopy(t *testing.T) {

	CopyInternal(testFile, "/home/aleksandr/Загрузки/vks-client02-ios.ovpn.bcp", 10, 1025)
}
