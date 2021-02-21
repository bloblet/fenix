package utils

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

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

func LoadConfig(fileName string) *Config {

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		log.WithFields(
			log.Fields{
				"fileName": fileName,
				"err":      err,
			},
		).Panic("Error opening config file")
	}

	defer f.Close()

	config := Config{
		Logger: Logger{
			LogLevel: "error",
		},
	}

	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	err = decoder.Decode(&config)

	if err != nil {
		log.WithFields(
			log.Fields{
				"err": err,
			},
		).Panic("Error parsing config file")
	}

	return &config
}
