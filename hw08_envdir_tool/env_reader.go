package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrDirIsNotLoad = errors.New("DirectoryIsNotLoad")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println("cant load dir")
		return nil, ErrDirIsNotLoad
	}
	result := make(Environment, len(entries))
	for _, file := range entries {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}
		fileName := path.Join(dir, file.Name())
		file, err := os.Open(fileName)
		if err != nil {
			log.Println("Ошибка при чтении файла\n %v", err)
		}
		scanner := bufio.NewScanner(file)
		ev := EnvValue{}
		ev.NeedRemove = scanner.Scan()
		if ev.NeedRemove {
			result[file.Name()] = ev
			continue
		}
		ev.Value = scanner.Text()
		file.Close()
		ev.Value = strings.TrimRight(ev.Value, " \t")
		ev.Value = strings.ReplaceAll(ev.Value, "\x00", "\n")
		ev.NeedRemove = len(ev.Value) == 0
		result[file.Name()] = ev
	}
	return result, nil
}

func SetEnv(envs Environment) {
	for key, env := range envs {
		if env.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				log.Fatalf("Не смог удалить %s", key)
			}
			continue
		}
		err := os.Setenv(key, env.Value)
		if err != nil {
			log.Fatalf("Не смог установить %s", key)
		}
	}
}
