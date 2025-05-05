package httpclient

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var client *http.Client

func init() {
	httpProxy := os.Getenv("HTTP_PROXY_URL")
	httpsProxy := os.Getenv("HTTPS_PROXY_URL")
	var proxyURL *url.URL
	var err error
	// Pick https proxy over http proxy
	if httpProxy != "" {
		proxyURL, err = url.Parse(httpProxy)
	} else if httpsProxy != "" {
		proxyURL, err = url.Parse(httpsProxy)
	}
	if err != nil {
		log.Println("failed to parse proxy URL: %w\n", err)
	}
	// Create http transport that uses the proxy
	transport := &http.Transport{}
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
		transport.MaxIdleConns = 100
		transport.MaxIdleConnsPerHost = 10
		transport.IdleConnTimeout = 90 * time.Second
		transport.DialContext = (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
	} else {
		transport.MaxIdleConns = 100
		transport.MaxIdleConnsPerHost = 10
		transport.IdleConnTimeout = 90 * time.Second
		transport.DialContext = (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
	}
	client = &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

// GetClient returns the singleton HTTP client
func GetClient() *http.Client {
	return client
}
