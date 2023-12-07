package helpers

import (
	"crypto/tls"
	"net/http"
	"time"
)

func InitHttpClient() *http.Client {
	transportCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   20 * time.Second,
		Transport: transportCfg,
	}
	return client
}
