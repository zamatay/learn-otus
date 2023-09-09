package main

import (
	"errors"
	"flag"
	"os"
	"path"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")

}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func main() {
	flag.Parse()
	if from == "" {
		to := path.Join(os.TempDir(), "tmmp.txt")
		Copy("./testData/input.txt", to, 0, 0)
	} else if FileExists(from) && FileExists(to) {
		Copy(from, to, offset, limit)
	}
}
