package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomain(unmarshalUser(readBuf(r)), domain)
}

type users [100_000]User

func readBuf(r io.Reader) <-chan string {
	//slog.Info("begin readBuf")
	out := make(chan string)

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			s := scanner.Text()
			//slog.String("scan string", s)
			out <- s
		}
		close(out)
	}()
	//slog.Info("end readBuf")
	return out
}

func unmarshalUser(in <-chan string) <-chan User {
	//slog.Info("begin unmarshalUser")
	out := make(chan User)
	go func() {
		for line := range in {
			var user User
			if err := json.Unmarshal([]byte(line), &user); err != nil {
				//slog.Error("error unmarshal", //slog.String("line", line))
				continue
			}
			//slog.Info("send User", //slog.Any("user", user))
			out <- user
		}
		close(out)
	}()
	//slog.Info("end unmarshalUser")
	return out
}

func countDomain(in <-chan User, domain string) (DomainStat, error) {
	//slog.Info("begin countDomain")
	result := make(DomainStat, 100_000)
	d := "" + domain
	for user := range in {
		matched := strings.Contains(user.Email, d)
		if matched {
			dom := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[dom]
			num++
			result[dom] = num
		}
	}
	//slog.Info("end countDomain")
	return result, nil
}
