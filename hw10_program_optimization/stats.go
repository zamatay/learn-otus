package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomain(unmarshalUser(readBuf(r)), domain)
}

func readBuf(r io.Reader) <-chan []byte {
	out := make(chan []byte)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			b := make([]byte, len(s.Bytes()))
			copy(b, s.Bytes())
			out <- b
			//out <- s.Bytes()
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
