package hw10programoptimization

import (
	"bufio"
	jsoniter "github.com/json-iterator/go"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomain(unmarshalUser(readBuf(r)), domain)
}

const workerCount = 2

func copyBuf(b []byte) []byte {
	result := make([]byte, len(b))
	copy(result, b)
	return result
}

func readBuf(r io.Reader) <-chan []byte {
	out := make(chan []byte, workerCount)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			out <- copyBuf(s.Bytes())
		}
		close(out)
	}()
	return out
}

func unmarshalUser(in <-chan []byte) <-chan User {
	out := make(chan User)
	go func() {
		var d []byte
		var user *User
		for d = range in {
			if err := jsoniter.Unmarshal(d, &user); err != nil {
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
		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(user.Email[strings.Index(user.Email, "@")+1:])]++
		}
	}
	return result, nil
}
