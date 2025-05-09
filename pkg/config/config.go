package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type configType struct {
	MongoDB struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"mongodb"`
	MetricsDB struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"metricsdb"`
}

var config configType

// GetMongodbHost returns configured mongodb host
func GetMongodbHost() string {
	return config.MongoDB.Host
}

// GetMongodbPort returns configured mongodb port
func GetMongodbPort() string {
	return config.MongoDB.Port
}

// GetMetricsdbHost returns configured metricsdb host
func GetMetricsdbHost() string {
	return config.MetricsDB.Host
}

// GetMetricsdbPort returns configured metricsdb port
func GetMetricsdbPort() string {
	return config.MetricsDB.Port
}

// validateConfigPath just makes sure, that the path provided is a file,
// that can be read
func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseConfig returns a new decoded Config struct
func ParseConfig(configPath string) error {
	// validate config path before decoding
	if err := validateConfigPath(configPath); err != nil {
		return err
	}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return err
	}

	return nil
}
