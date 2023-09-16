package main

import (
	"errors"
	"log"
	"os"
)

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Не хватает аргументов")
	}
	pathToDir := os.Args[1]
	if pathToDir == "" {
		log.Fatalf("Path is empty")
	}
	if !FileExists(pathToDir) {
		log.Fatal("Dir not exists or not rights")
		return
	}
	env, _ := ReadDir(pathToDir)
	RunCmd(os.Args[2:], env)
}
