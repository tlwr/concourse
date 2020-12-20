package http

import (
	"fmt"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/concourse/atc/creds"
)

type HTTPManager struct {
	URL string `long:"url" description:"HTTP server address used to access secrets."`
}

func (manager *HTTPManager) Init(log lager.Logger) error {
	return nil
}

func (manager HTTPManager) IsConfigured() bool {
	return manager.URL != ""
}

func (manager HTTPManager) Validate() error {
	parsedUrl, err := url.Parse(manager.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", err)
	}

	// "foo" will parse without error (as a Path, with an empty Host)
	// so we'll do a few additional sanity checks that this is a valid URL
	// copied from atc/credhub/manager.go
	if parsedUrl.Host == "" || !(parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https") {
		return fmt.Errorf("invalid URL (must be http or https)")
	}

	return nil
}

func (manager HTTPManager) Health() (*creds.HealthResponse, error) {
	path := "/health"

	healthResponse := &creds.HealthResponse{
		Method: path,
	}

	resp, err := http.Get(manager.URL + path)
	if err != nil {
		return healthResponse, err
	}
	defer resp.Body.Close()

	healthResponse.Response = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		return healthResponse, fmt.Errorf("expected HTTP 200 received %d", resp.StatusCode)
	}

	return healthResponse, nil
}

func (manager HTTPManager) NewSecretsFactory(logger lager.Logger) (creds.SecretsFactory, error) {
	return NewHTTPFactory(logger, manager.URL), nil
}

func (manager HTTPManager) Close(logger lager.Logger) {
	// TODO - to implement
}
