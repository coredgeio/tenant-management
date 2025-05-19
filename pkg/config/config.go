package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type DefaultHeaders struct {
	Authorization string `yaml:"authorization"`
	ContentType   string `yaml:"contentType"`
	Apikey        string `yaml:"apikey"`
}

type EndpointDetails struct {
	BaseUrl        string         `yaml:"baseUrl"`
	HttpMethod     string         `yaml:"httpMethod"`
	DefaultHeaders DefaultHeaders `yaml:"defaultHeaders"`
}

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
	PollingTime    int `yaml:"pollingTime"`
	TenantMetadata struct {
		Enabled bool `yaml:"enabled"`
		// this field makes sure that we stop polling once we have received our desired response if set to true
		StopUpdateOnceSet bool            `yaml:"stopUpdateOnceSet"`
		EndpointDetails   EndpointDetails `yaml:"endpointDetails"`
	} `yaml:"tenantMetadata"`
	TenantUserMetadata struct {
		Enabled bool `yaml:"enabled"`
		// this field makes sure that we stop polling once we have received our desired response if set to true
		StopUpdateOnceSet bool            `yaml:"stopUpdateOnceSet"`
		EndpointDetails   EndpointDetails `yaml:"endpointDetails"`
	} `yaml:"tenantUserMetadata"`
	PublishMeteringInfo struct {
		Enabled         bool            `yaml:"enabled"`
		EndpointDetails EndpointDetails `yaml:"endpointDetails"`
	} `yaml:"publishMeteringInfo"`
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

// GetPollingTime returns configured polling time
func GetPollingTime() int {
	return config.PollingTime
}

// GetTenantMetadataEnabled returns if tenant metadata is enabled
func GetTenantMetadataEnabled() bool {
	return config.TenantMetadata.Enabled
}

// GetTenantMetadataStopUpdateOnceSet returns if tenant metadata update should stop once set
func GetTenantMetadataStopUpdateOnceSet() bool {
	return config.TenantMetadata.StopUpdateOnceSet
}

// GetTenantMetadataEndpointDetails returns tenant metadata endpoint details
func GetTenantMetadataEndpointDetails() EndpointDetails {
	return config.TenantMetadata.EndpointDetails
}

// GetTenantUserMetadataEnabled returns if tenant user metadata is enabled
func GetTenantUserMetadataEnabled() bool {
	return config.TenantUserMetadata.Enabled
}

// GetTenantUserMetadataStopUpdateOnceSet returns if tenant user metadata update should stop once set
func GetTenantUserMetadataStopUpdateOnceSet() bool {
	return config.TenantUserMetadata.StopUpdateOnceSet
}

// GetTenantUserMetadataEndpointDetails returns tenant user metadata endpoint details
func GetTenantUserMetadataEndpointDetails() EndpointDetails {
	return config.TenantUserMetadata.EndpointDetails
}

// GetPublishMeteringInfoEnabled returns if publish metering info is enabled
func GetPublishMeteringInfoEnabled() bool {
	return config.PublishMeteringInfo.Enabled
}

// GetPublishMeteringInfoEndpointDetails returns publish metering info endpoint details
func GetPublishMeteringInfoEndpointDetails() EndpointDetails {
	return config.PublishMeteringInfo.EndpointDetails
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
