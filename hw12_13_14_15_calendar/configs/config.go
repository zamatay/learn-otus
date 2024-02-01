package configs

import (
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var cfg = &Config{}
var once = sync.Once{}

type HTTP struct {
	Host string `yaml:"host" env-default:"http://localhost"`
	Port string `yaml:"port"`
}

type Broker struct {
	Url      string `yaml:"url"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

type Schedule struct {
	Interval time.Duration `yaml:"interval"`
}

type Config struct {
	DevEnv   string       `yaml:"env" env-default:"local"`
	Log      LoggerConfig `yaml:"logger"`
	DB       DBConfig     `yaml:"db"`
	HTTP     HTTP         `yaml:"http"`
	Grpc     Grpc         `yaml:"grpc"`
	Broker   Broker       `yaml:"broker"`
	Schedule Schedule     `yaml:"schedule"`
}

type Grpc struct {
	Port int
}

type LoggerConfig struct {
	Level string `yaml:"level" env-default:"string"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	AppName  string `yaml:"appName"`
}

const defaultConfig = "./configs/config.yaml"

func Configs() *Config {
	once.Do(func() {
		cfg = NewConfig()
	})
	return cfg
}

func NewConfig() *Config {
	configPath, err := fetchConfigPath()
	if err != nil {
		return &Config{}
	}
	log.Print(os.Getwd())
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

func fetchConfigPath() (string, error) {
	var res string
	for i := 0; i < len(os.Args); i++ {
		if strings.HasPrefix(strings.ToLower(os.Args[i]), "-config") {
			t := strings.Split(os.Args[i], "=")
			if len(t) > 1 {
				res = t[1]
				return res, nil
			}
		}
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" && fileExists(defaultConfig) {
		res = defaultConfig
	}

	if res == "" {
		return "", errors.New("path not define")
	}

	return res, nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
