// Package http_helper contains helpers to interact with deployed resources through HTTP.
package http_helper

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
)

// HttpGet performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code and body. If there's any error, fail the test.
func HttpGet(t *testing.T, url string, tlsConfig *tls.Config) (int, string) {
	statusCode, body, err := HttpGetE(t, url, tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
	return statusCode, body
}

// HttpGetE performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code, body, and any error.
func HttpGetE(t *testing.T, url string, tlsConfig *tls.Config) (int, string, error) {
	logger.Logf(t, "Making an HTTP GET call to URL %s", url)

	// Set HTTP client transport config
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout: 10 * time.Second,
		// Include the previously created transport config
		Transport: tr,
	}

	resp, err := client.Get(url)
	if err != nil {
		return -1, "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, strings.TrimSpace(string(body)), nil
}

// HttpGetWithValidation performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, fail the test.
func HttpGetWithValidation(t *testing.T, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) {
	err := HttpGetWithValidationE(t, url, tlsConfig, expectedStatusCode, expectedBody)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithValidationE performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, return an error.
func HttpGetWithValidationE(t *testing.T, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) error {
	return HttpGetWithCustomValidationE(t, url, tlsConfig, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	})
}

// HttpGetWithCustomValidation performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidation(t *testing.T, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) {
	err := HttpGetWithCustomValidationE(t, url, tlsConfig, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithCustomValidationE performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidationE(t *testing.T, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) error {
	statusCode, body, err := HttpGetE(t, url, tlsConfig)

	if err != nil {
		return err
	}

	if !validateResponse(statusCode, body) {
		return ValidationFunctionFailed{Url: url, Status: statusCode, Body: body}
	}

	return nil
}

// HttpGetWithRetry repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetry(t *testing.T, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	err := HttpGetWithRetryE(t, url, tlsConfig, expectedStatus, expectedBody, retries, sleepBetweenRetries)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryE repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetryE(t *testing.T, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithValidationE(t, url, tlsConfig, expectedStatus, expectedBody)
	})

	return err
}

// HttpGetWithRetryWithCustomValidation repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidation(t *testing.T, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	err := HttpGetWithRetryWithCustomValidationE(t, url, tlsConfig, retries, sleepBetweenRetries, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryWithCustomValidationE repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidationE(t *testing.T, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithCustomValidationE(t, url, tlsConfig, validateResponse)
	})

	return err
}
