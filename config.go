package main

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Config the configuration
type Config struct {
	// Both the server configuration and the endpoint configuration for the client
	Server struct {
		Port    int
		Address string
	}
	// File system configuration
	Fs struct {
		// The root being updated for server mode, or monitored for client mode
		RootDir string `yaml:"root_dir"`
		// The path to the password file
		PasswordFile string `yaml:"password_file"`
	}
}

// LoadConfig loads the config
func LoadConfig() *Config {
	data, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal("Unable to load config file")
		os.Exit(1)
	}
	var config Config
	err2 := yaml.Unmarshal(data, &config)
	if err2 != nil {
		log.Fatal("Failed to parse config file")
		os.Exit(1)
	}
	return &config
}
