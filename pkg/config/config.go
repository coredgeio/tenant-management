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
	Client struct {
		Name string `yaml:"name"`
	} `yaml:"client"`
	TenantLevelKYC struct {
		Enabled           bool   `yaml:"enabled"`
		ServerPath        string `yaml:"serverPath"`
		PollingTime       int    `yaml:"pollingTime"`
		StopUpdateOnceSet bool   `yaml:"stopUpdateOnceSet"`
		RequestDetails    struct {
			BaseUrl        string `yaml:"baseUrl"`
			HttpMethod     string `yaml:"httpMethod"`
			DefaultHeaders struct {
				Authorization string `yaml:"authorization"`
				ContentType   string `yaml:"contentType"`
			} `yaml:"defaultHeaders"`
		} `yaml:"requestDetails"`
	} `yaml:"tenantLevelKyc"`
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

// GetClientName returns configured client name
func GetClientName() string {
	return config.Client.Name
}

// GetTenantLevelKYCEnabled returns configured TenantLevelKYC enabled
func GetTenantLevelKYCEnabled() bool {
	return config.TenantLevelKYC.Enabled
}

// GetTenantLevelKYCServerPath returns configured TenantLevelKYC server path
func GetTenantLevelKYCServerPath() string {
	return config.TenantLevelKYC.ServerPath
}

// GetTenantLevelKYCPollingTime returns configured TenantLevelKYC polling time
func GetTenantLevelKYCPollingTime() int {
	return config.TenantLevelKYC.PollingTime
}

// GetTenantLevelKYCStopUpdateOnceSet returns configured TenantLevelKYC stop update once set
func GetTenantLevelKYCStopUpdateOnceSet() bool {
	return config.TenantLevelKYC.StopUpdateOnceSet
}

// GetBaseUrl returns the base URL
func GetTenantLevelKYCBaseUrl() string {
	return config.TenantLevelKYC.RequestDetails.BaseUrl
}

// GetHttpMethod returns the HTTP method (GET, POST, etc.)
func GetTenantLevelKYCHttpMethod() string {
	return config.TenantLevelKYC.RequestDetails.HttpMethod
}

// GetAuthorization returns the Authorization header value
func GetTenantLevelKYCAuthorization() string {
	return config.TenantLevelKYC.RequestDetails.DefaultHeaders.Authorization
}

// GetContentType returns the Content-Type header value
func GetTenantLevelKYCContentType() string {
	return config.TenantLevelKYC.RequestDetails.DefaultHeaders.ContentType
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
