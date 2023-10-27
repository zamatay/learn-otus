package hw10programoptimization

import (
	"bufio"
	"github.com/mailru/easyjson"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	s := bufio.NewScanner(r)
	var user User
	for s.Scan() {
		if err := easyjson.Unmarshal(s.Bytes(), &user); err != nil {
			continue
		}
		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(user.Email[strings.Index(user.Email, "@")+1:])]++
		}
	}
	return result, nil
}
