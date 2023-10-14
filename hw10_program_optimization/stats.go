package hw10programoptimization

import (
	"bufio"
	jsoniter "github.com/json-iterator/go"
	"io"
	"strings"
	"sync"
)

type User struct {
	Email string
}

var wg sync.WaitGroup

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomain(unmarshalUser(readBuf(r)), domain)
}

func readBuf(r io.Reader) <-chan []byte {
	out := make(chan []byte)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			wg.Add(1)
			out <- s.Bytes()
			wg.Wait()
		}
		close(out)
	}()
	return out
}

func unmarshalUser(in <-chan []byte) <-chan User {
	out := make(chan User)
	go func() {
		var user *User
		var line []byte
		for line = range in {
			if err := jsoniter.Unmarshal(line, &user); err != nil {
				continue
			}
			wg.Done()
			out <- *user
		}
		close(out)
	}()
	return out
}

func countDomain(in <-chan User, domain string) (DomainStat, error) {
	result := make(DomainStat)
	var user User
	for user = range in {
		if strings.Contains(user.Email, domain) {
			result[strings.ToLower(user.Email[strings.Index(user.Email, "@")+1:])]++
		}
	}
	return result, nil
}
