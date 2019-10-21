package taxis99

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestCostCenterFind(t *testing.T) {
	testPath(t, string(costCentersEndpoint), func(c *Client) error {
		_, err := c.CostCenter.Find(context.Background(), nil)
		return err
	})

	testMethod(t, http.MethodGet, func(c *Client) error {
		_, err := c.CostCenter.Find(context.Background(), nil)
		return err
	})

	testQuery(t, []Filter{
		{"search": "123"},
		{"search": "124", "limit": "100"},
	}, ccFields, func(c *Client, f Filter) error {
		_, err := c.CostCenter.Find(context.Background(), f)
		return err
	})

	testResponseBody(t, [][]byte{
		[]byte(`[{"id":123,"name":"TI","enabled":true,"company":{"id":"1234","name":"Mobilitee"}}]`),
	}, func(c *Client) (interface{}, error) {
		return c.CostCenter.Find(context.Background(), nil)
	})
}

func TestCostCenterFindError(t *testing.T) {
	testError(t, func(c *Client) error {
		_, err := c.CostCenter.Find(context.Background(), nil)
		return err
	})
}

func TestCostCenterCreate(t *testing.T) {
	testPath(t, string(costCentersEndpoint), func(c *Client) error {
		_, err := c.CostCenter.Create(context.Background(), CostCenter{})
		return err
	})

	testMethod(t, http.MethodPost, func(c *Client) error {
		_, err := c.CostCenter.Create(context.Background(), CostCenter{})
		return err
	})

	testResponseBody(t, [][]byte{
		[]byte(`{"id":123,"name":"IT","enabled":true,"company":{"id":"1234","name":"Mobilitee"}}`),
		[]byte(`{"id":124,"name":"Marketing","enabled":true,"company":{"id":"1234","name":"Mobilitee"}}`),
	}, func(c *Client) (interface{}, error) {
		return c.CostCenter.Create(context.Background(), CostCenter{})
	})

	testRequestBody(t, []func(*Client) ([]byte, error){
		func(c *Client) (want []byte, err error) {
			want = []byte(`{"name":"IT"}`)
			_, err = c.CostCenter.Create(context.Background(), CostCenter{Name: "IT"})
			return
		},
		func(c *Client) (want []byte, err error) {
			want = []byte(`{"name":"Sales"}`)
			_, err = c.CostCenter.Create(context.Background(), CostCenter{Name: "Sales"})
			return
		},
	})
}

func TestCostCenterCreateError(t *testing.T) {
	testError(t, func(c *Client) error {
		_, err := c.CostCenter.Create(context.Background(), CostCenter{})
		return err
	})
}

func TestCostCenterRemove(t *testing.T) {
	testCases := []struct {
		id   int64
		want string
	}{
		{25, fmt.Sprintf(string(costCenterEndpoint), 25)},
		{28, fmt.Sprintf(string(costCenterEndpoint), 28)},
	}

	for _, tc := range testCases {
		testPath(t, tc.want, func(c *Client) error {
			return c.CostCenter.Remove(context.Background(), tc.id)
		})
	}

	testMethod(t, http.MethodDelete, func(c *Client) error {
		return c.CostCenter.Remove(context.Background(), 20)
	})
}

func TestCostCenterRemoveError(t *testing.T) {
	testError(t, func(c *Client) error {
		return c.CostCenter.Remove(context.Background(), 0)
	})
}
