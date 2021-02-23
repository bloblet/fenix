package utils

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

var configOnce = sync.Once{}
var config *Config

type API struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Logger struct {
	LogLevel string `yaml:"loglevel"`
}

type Database struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

type Config struct {

	// API specific settings
	API API `yaml:"api"`

	Logger Logger `yaml:"logger"`

	// Database specific settings
	Database Database `yaml:"database"`
}

func LoadConfig() *Config {
	configOnce.Do(readConfig)

	return config
}

func readConfig() {

	f, err := os.OpenFile("fenix.yml", os.O_RDONLY, 0)

	if err != nil {
		log.WithFields(
			log.Fields{
				"fileName": "fenix.yml",
				"err":      err,
			},
		).Info("Error opening config file, using defaults.")

		config = &Config{
			API: API{
				Host: "localhost",
				Port: 4545,
			},
			Database: Database{
				Database: "development",
				Host: "localhost",
			},
			Logger: Logger{
				LogLevel: "error",
			},
		}
		return
	}

	defer f.Close()

	c := Config{
		Logger: Logger{
			LogLevel: "error",
		},
	}

	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	err = decoder.Decode(&c)

	if err != nil {
		log.WithFields(
			log.Fields{
				"err": err,
			},
		).Panic("Error parsing config file")
	}

	config = &c
}
