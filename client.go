package taxis99

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.corp.99taxis.com/v2/"
)

// ApiError implements the error interface
// and returns infos from the request.
type APIError struct {
	StatusCode int
	Msg        string
	Err        error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Error Status Code: %d; Message: %s.", e.StatusCode, e.Msg)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// requester is the interface that performs a request
// to the server and delegates the parsing to the parser interface.
type requester interface {
	Request(ctx context.Context, method, path string, body, output interface{}) error
}

type service struct {
	client requester
}

// Client is responsible for handling request to the Taxis 99 API.
type Client struct {
	// client to connect to the API.
	client *http.Client

	// Host used for API requests.
	// Should always be specified with a trailing slash.
	BaseURL *url.URL

	// reuse a single struct instead of allocating one for each service on the heap.
	common service

	CostCenter *CostCenterService
}

// NewClient returns a reference to the Client struct.
func NewClient(hc *http.Client) *Client {
	if hc == nil {
		hc = http.DefaultClient
	}
	// Deep copy the the URL.
	u, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:  hc,
		BaseURL: u,
	}

	c.common.client = c

	c.CostCenter = (*CostCenterService)(&c.common)

	return c
}

// Request created an API request. A relative path can be providaded
// in which case it is resolved relative to the host of the Client.
func (c *Client) Request(ctx context.Context, method, path string, body, output interface{}) error {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), &buf)
	if err != nil {
		return err
	}

	if buf.Len() > 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Ignores io.EOF error caused by empty response body.
	if err = json.NewDecoder(res.Body).Decode(output); err != nil && !errors.Is(err, io.EOF) {
		return &APIError{
			StatusCode: res.StatusCode,
			Msg:        fmt.Sprintf("api: '%s'.", err.Error()),
			Err:        err,
		}
	}

	return nil
}
