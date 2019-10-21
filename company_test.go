package taxis99

import (
	"context"
	"net/http"
	"testing"
)

func TestCompanyFind(t *testing.T) {
	testPath(t, string(companiesEndpoint), func(c *Client) error {
		_, err := c.Company.Find(context.Background())
		return err
	})

	testMethod(t, http.MethodGet, func(c *Client) error {
		_, err := c.Company.Find(context.Background())
		return err
	})

	testResponseBody(t, [][]byte{
		[]byte(`[{"id":"123","name":"Mobilitee"}]`),
	}, func(c *Client) (interface{}, error) {
		return c.Company.Find(context.Background())
	})
}

func TestCompanyFindError(t *testing.T) {
	testError(t, func(c *Client) error {
		_, err := c.Company.Find(context.Background())
		return err
	})
}
