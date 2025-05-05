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
				Apikey        string `yaml:"apikey"`
			} `yaml:"defaultHeaders"`
		} `yaml:"requestDetails"`
	} `yaml:"tenantLevelKyc"`
	PaymentConfigurationForTenant struct {
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
				Apikey        string `yaml:"apikey"`
			} `yaml:"defaultHeaders"`
		} `yaml:"requestDetails"`
	} `yaml:"paymentConfigurationForTenant"`
	TenantType struct {
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
				Apikey        string `yaml:"apikey"`
			} `yaml:"defaultHeaders"`
		} `yaml:"requestDetails"`
	} `yaml:"tenantType"`
	TenantUserLevelKYC struct {
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
				Apikey        string `yaml:"apikey"`
			} `yaml:"defaultHeaders"`
		} `yaml:"requestDetails"`
	} `yaml:"tenantUserLevelKyc"`
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

// GetTenantLevelKYCApiKey returns the apiKey header value
func GetTenantLevelKYCApiKey() string {
	return config.TenantLevelKYC.RequestDetails.DefaultHeaders.Apikey
}

// GetPaymentMethodConfigurationEnabled returns configured Payment Configuration enabled
func GetPaymentMethodConfigurationEnabled() bool {
	return config.PaymentConfigurationForTenant.Enabled
}

// GetPaymentMethodConfigurationServerPath returns configured payment configuration server path
func GetPaymentMethodConfigurationServerPath() string {
	return config.PaymentConfigurationForTenant.ServerPath
}

// GetPaymentMethodConfigurationPollingTime returns configured Payment configuration polling time
func GetPaymentMethodConfigurationPollingTime() int {
	return config.PaymentConfigurationForTenant.PollingTime
}

// GetPaymentMethodConfigurationStopUpdateOnceSet returns configured Payment Configuration stop update once set
func GetPaymentMethodConfigurationStopUpdateOnceSet() bool {
	return config.PaymentConfigurationForTenant.StopUpdateOnceSet
}

// GetPaymentMethodConfigurationBaseUrl returns the base URL
func GetPaymentMethodConfigurationBaseUrl() string {
	return config.PaymentConfigurationForTenant.RequestDetails.BaseUrl
}

// GetPaymentMethodConfigurationHttpMethod returns the HTTP method (GET, POST, etc.)
func GetPaymentMethodConfigurationHttpMethod() string {
	return config.PaymentConfigurationForTenant.RequestDetails.HttpMethod
}

// GetPaymentMethodConfigurationAuthorization returns the Authorization header value
func GetPaymentMethodConfigurationAuthorization() string {
	return config.PaymentConfigurationForTenant.RequestDetails.DefaultHeaders.Authorization
}

// GetPaymentMethodConfigurationContentType returns the Content-Type header value
func GetPaymentMethodConfigurationContentType() string {
	return config.PaymentConfigurationForTenant.RequestDetails.DefaultHeaders.ContentType
}

// GetPaymentMethodConfigurationApiKey returns the apiKey header value
func GetPaymentMethodConfigurationApiKey() string {
	return config.PaymentConfigurationForTenant.RequestDetails.DefaultHeaders.Apikey
}

// GetTenantTypeEnabled returns configured Payment Configuration enabled
func GetTenantTypeEnabled() bool {
	return config.TenantType.Enabled
}

// GetTenantTypeServerPath returns configured payment configuration server path
func GetTenantTypeServerPath() string {
	return config.TenantType.ServerPath
}

// GetTenantTypePollingTime returns configured Payment configuration polling time
func GetTenantTypePollingTime() int {
	return config.TenantType.PollingTime
}

// GetTenantTypeStopUpdateOnceSet returns configured Payment Configuration stop update once set
func GetTenantTypeStopUpdateOnceSet() bool {
	return config.TenantType.StopUpdateOnceSet
}

// GetTenantTypeBaseUrl returns the base URL
func GetTenantTypeBaseUrl() string {
	return config.TenantType.RequestDetails.BaseUrl
}

// GetTenantTypeHttpMethod returns the HTTP method (GET, POST, etc.)
func GetTenantTypeHttpMethod() string {
	return config.TenantType.RequestDetails.HttpMethod
}

// GetTenantTypeAuthorization returns the Authorization header value
func GetTenantTypeAuthorization() string {
	return config.TenantType.RequestDetails.DefaultHeaders.Authorization
}

// GetTenantTypeContentType returns the Content-Type header value
func GetTenantTypeContentType() string {
	return config.TenantType.RequestDetails.DefaultHeaders.ContentType
}

// GetTenantTypeApiKey returns the apiKey header value
func GetTenantTypeApiKey() string {
	return config.TenantType.RequestDetails.DefaultHeaders.Apikey
}

// GetTenantUserLevelKycEnabled returns configured Payment Configuration enabled
func GetTenantUserLevelKYCEnabled() bool {
	return config.TenantUserLevelKYC.Enabled
}

// GetTenantUserLevelKYCServerPath returns configured payment configuration server path
func GetTenantUserLevelKYCServerPath() string {
	return config.TenantUserLevelKYC.ServerPath
}

// GetTenantUserLevelKYCPollingTime returns configured Payment configuration polling time
func GetTenantUserLevelKYCPollingTime() int {
	return config.TenantUserLevelKYC.PollingTime
}

// GetTenantUserLevelKYCStopUpdateOnceSet returns configured Payment Configuration stop update once set
func GetTenantUserLevelKYCStopUpdateOnceSet() bool {
	return config.TenantUserLevelKYC.StopUpdateOnceSet
}

// GetTenantUserLevelKYCBaseUrl returns the base URL
func GetTenantUserLevelKYCBaseUrl() string {
	return config.TenantUserLevelKYC.RequestDetails.BaseUrl
}

// GetTenantUserLevelKYCHttpMethod returns the HTTP method (GET, POST, etc.)
func GetTenantUserLevelKYCHttpMethod() string {
	return config.TenantUserLevelKYC.RequestDetails.HttpMethod
}

// GetTenantUserLevelKYCAuthorization returns the Authorization header value
func GetTenantUserLevelKYCAuthorization() string {
	return config.TenantUserLevelKYC.RequestDetails.DefaultHeaders.Authorization
}

// GetTenantUserLevelKYCContentType returns the Content-Type header value
func GetTenantUserLevelKYCContentType() string {
	return config.TenantUserLevelKYC.RequestDetails.DefaultHeaders.ContentType
}

// GetTenantUserLevelKYCApiKey returns the apiKey header value
func GetTenantUserLevelKYCApiKey() string {
	return config.TenantUserLevelKYC.RequestDetails.DefaultHeaders.Apikey
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
