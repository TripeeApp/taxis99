package taxis99

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// reqStatus is the request status
// Success = 1
// Fail = 0
type reqStatus int

const (
	ReqStatusFail reqStatus = iota
	ReqStatusOK

	// Error message format
	errFmt = `Error in request to 99taxis API: %s; HTTP Status code: %d; Body: %s;`
)

// the request status
type status struct {
	Status reqStatus `json:"status"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf(errFmt, e.msg, e.statusCode, e.body)
}

// APIError implement the commin error interface
type APIError struct {
	statusCode 	int
	body		[]byte
	msg			string
}

// APIConfiguration this struct can be used to pass in arguments to the constructor
// this exists so we can keep the arguments count of the new function down to 3
type APIConfiguration struct {
	APIKey 		string
	CompanyID 	string
}

type requester interface {
	Request(ctx context.Context, method string, path endpoint, body, output interface{}) error
}

// Client the client of the API
type Client struct {
	host *url.URL
	client *http.Client
	User *UserService
}

// New create a new Client instance
func New(host *url.URL, configuration APIConfiguration, c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}
	c.Transport = &Transport {
		configuration.APIKey,
		configuration.CompanyID,
		c.Transport,
	}

	client := &Client{host: host, client: c}
	client.User = &UserService{client}
	return client
}

// Request 
func (c *Client) Request(ctx context.Context, method string, path endpoint, body, output interface{}) error {
	u, err := c.host.Parse(path.String(ctx))
	if err != nil {
		return err
	}

	var b io.ReadWriter
	if body != nil {
		b = new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, u.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	r, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(r, output); err != nil {
		return &APIError{
			statusCode: res.StatusCode,
			body: r,
			msg: err.Error(),
		}
	}

	return nil
}