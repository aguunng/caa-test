package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	HTTPClient *http.Client
}

func defaultHTTPClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Logger = nil
	retryClient.HTTPClient.Timeout = 20 * time.Second
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	return retryClient.StandardClient()
}

func New() *Client {
	return &Client{
		HTTPClient: defaultHTTPClient(),
	}
}

func (c *Client) Call(ctx context.Context, method, url string, body io.Reader, headers map[string]string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), url, body)
	if err != nil {
		return &Error{
			Message:  fmt.Sprintf("unable to create new request: %s", err.Error()),
			RawError: err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return &Error{
			Message:  fmt.Sprintf("unable to sends an http request: %s", err.Error()),
			RawError: err,
		}
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Error{
			Message:    fmt.Sprintf("unable to read response body: %s", err.Error()),
			StatusCode: resp.StatusCode,
			RawError:   err,
		}
	}

	if resp.StatusCode >= 400 {
		return &Error{
			Message:    "http client error",
			StatusCode: resp.StatusCode,
			RawError: fmt.Errorf(
				"%s %s returned error %d response: %s",
				resp.Request.Method,
				resp.Request.URL,
				resp.StatusCode,
				string(responseBody)),
		}
	}

	if response != nil {
		if err = json.Unmarshal(responseBody, response); err != nil {
			return &Error{
				Message:    fmt.Sprintf("unable to unmarshaling body response: %s", err.Error()),
				StatusCode: resp.StatusCode,
				RawError:   err,
			}
		}
	}

	return nil
}
