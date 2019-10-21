package taxis99

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestEmployeeFind(t *testing.T) {
	testPath(t, string(employeesEndpoint), func(c *Client) error {
		_, err := c.Employee.Find(context.Background(), nil)
		return err
	})

	testMethod(t, http.MethodGet, func(c *Client) error {
		_, err := c.Employee.Find(context.Background(), nil)
		return err
	})

	testQuery(t, []Filter{
		{"search": "123"},
		{"search": "124", "limit": "100"},
		{"search": "124", "limit": "100", "invalid": "param"},
	}, employeeFields, func(c *Client, f Filter) error {
		_, err := c.Employee.Find(context.Background(), f)
		return err
	})

	testResponseBody(t, [][]byte{
		[]byte(`[{"id":125,"name":"José Santos","email":"jose.santos@empresa.com.br","phone":{"number":"11999999999","country":"BRA"},"company":{"id":"47a3083b-5d03-4e05-ad9d-9fd6fddd613e","name":"99"},"nationalId":"98765432100","supervisorId":167,"enabled":true,"externalId":0,"categories":["regular-taxi","top99","turbo-taxi","pop99"]}]`),
	}, func(c *Client) (interface{}, error) {
		return c.Employee.Find(context.Background(), nil)
	})
}

func TestEmployeeFindError(t *testing.T) {
	testError(t, func(c *Client) error {
		_, err := c.Employee.Find(context.Background(), nil)
		return err
	})
}

func TestEmployeeCreate(t *testing.T) {
	testPath(t, string(employeesEndpoint), func(c *Client) error {
		_, err := c.Employee.Create(context.Background(), Employee{}, false)
		return err
	})

	testMethod(t, http.MethodPost, func(c *Client) error {
		_, err := c.Employee.Create(context.Background(), Employee{}, false)
		return err
	})

	testResponseBody(t, [][]byte{
		[]byte(`{"employee":{"name":"José Santos","email":"jose.santos@empresa.com.br","phone":{"number":"11999999999","country":"BRA"},"nationalId":"98765432100","externalId":55091,"categories":["regular-taxi","turbo-taxi","pop99"]},"sendWelcomeEmail":false}`),
	}, func(c *Client) (interface{}, error) {
		return c.Employee.Create(context.Background(), Employee{}, false)
	})

	testRequestBody(t, []func(*Client) ([]byte, error){
		func(c *Client) (want []byte, err error) {
			want = []byte(`{"employee":{"name":"José Santos","email":"jose.santos@empresa.com.br","phone":{"number":"11999999999","country":"BRA"},"nationalId":"98765432100","externalId":55091,"categories":["regular-taxi","turbo-taxi","pop99"]},"sendWelcomeEmail":false}`)
			_, err = c.Employee.Create(context.Background(), Employee{
				Name:  "José Santos",
				Email: "jose.santos@empresa.com.br",
				Phone: &Phone{
					Number:  "11999999999",
					Country: "BRA",
				},
				NationalID: "98765432100",
				ExternalID: 55091,
				Categories: []string{"regular-taxi", "turbo-taxi", "pop99"},
			}, false)
			return
		},
		func(c *Client) (want []byte, err error) {
			want = []byte(`{"employee":{"name":"José Santos","email":"jose.santos@empresa.com.br","phone":{"number":"11999999999","country":"BRA"},"nationalId":"98765432100","externalId":55092,"categories":["regular-taxi","turbo-taxi","pop99"]},"sendWelcomeEmail":true}`)
			_, err = c.Employee.Create(context.Background(), Employee{
				Name:  "José Santos",
				Email: "jose.santos@empresa.com.br",
				Phone: &Phone{
					Number:  "11999999999",
					Country: "BRA",
				},
				NationalID: "98765432100",
				ExternalID: 55092,
				Categories: []string{"regular-taxi", "turbo-taxi", "pop99"},
			}, true)
			return
		},
	})
}

func TestEmployeeCreateError(t *testing.T) {
	testError(t, func(c *Client) error {
		_, err := c.Employee.Create(context.Background(), Employee{}, false)
		return err
	})
}

func TestEmployeeUpdate(t *testing.T) {
	testCases := []struct {
		id   int64
		want string
	}{
		{25, fmt.Sprintf(string(employeeEndpoint), 25)},
		{28, fmt.Sprintf(string(employeeEndpoint), 28)},
	}

	for _, tc := range testCases {
		testPath(t, tc.want, func(c *Client) error {
			_, err := c.Employee.Update(context.Background(), Employee{ID: tc.id})
			return err
		})
	}

	testMethod(t, http.MethodPut, func(c *Client) error {
		_, err := c.Employee.Update(context.Background(), Employee{})
		return err
	})

	testResponseBody(t, [][]byte{
		[]byte(`{"name":"José Santos","email":"jose.santos@empresa.com.br","phone":{"number":"11999999999","country":"BRA"},"nationalId":"98765432100","externalId":55091,"categories":["regular-taxi","turbo-taxi","pop99"]}`),
	}, func(c *Client) (interface{}, error) {
		return c.Employee.Update(context.Background(), Employee{})
	})

	testRequestBody(t, []func(*Client) ([]byte, error){
		func(c *Client) (want []byte, err error) {
			want = []byte(`{"employee":{"id":10,"name":"José Santos","email":"jose.santos@empresa.com.br","phone":{"number":"11999999999","country":"BRA"},"nationalId":"98765432100","externalId":55091,"categories":["regular-taxi","turbo-taxi","pop99"]},"sendWelcomeEmail":false}`)
			_, err = c.Employee.Update(context.Background(), Employee{
				ID:    10,
				Name:  "José Santos",
				Email: "jose.santos@empresa.com.br",
				Phone: &Phone{
					Number:  "11999999999",
					Country: "BRA",
				},
				NationalID: "98765432100",
				ExternalID: 55091,
				Categories: []string{"regular-taxi", "turbo-taxi", "pop99"},
			})
			return
		},
		func(c *Client) (want []byte, err error) {
			want = []byte(`{"employee":{"id":100,"name":"José Santos","email":"jose.santos@empresa.com.br","phone":{"number":"11999999999","country":"BRA"},"nationalId":"98765432100","externalId":55091,"categories":["regular-taxi","turbo-taxi","pop99"]},"sendWelcomeEmail":false}`)
			_, err = c.Employee.Update(context.Background(), Employee{
				ID:    100,
				Name:  "José Santos",
				Email: "jose.santos@empresa.com.br",
				Phone: &Phone{
					Number:  "11999999999",
					Country: "BRA",
				},
				NationalID: "98765432100",
				ExternalID: 55091,
				Categories: []string{"regular-taxi", "turbo-taxi", "pop99"},
			})
			return
		},
	})
}

func TestEmployeeUpdateError(t *testing.T) {
	testError(t, func(c *Client) error {
		_, err := c.Employee.Update(context.Background(), Employee{})
		return err
	})
}

func TestEmployeeRemove(t *testing.T) {
	testCases := []struct {
		id   int64
		want string
	}{
		{25, fmt.Sprintf(string(employeeEndpoint), 25)},
		{28, fmt.Sprintf(string(employeeEndpoint), 28)},
	}

	for _, tc := range testCases {
		testPath(t, tc.want, func(c *Client) error {
			return c.Employee.Remove(context.Background(), tc.id)
		})
	}

	testMethod(t, http.MethodDelete, func(c *Client) error {
		return c.Employee.Remove(context.Background(), 20)
	})
}

func TestEmployeeRemoveError(t *testing.T) {
	testError(t, func(c *Client) error {
		return c.Employee.Remove(context.Background(), 0)
	})
}

func TestEmployeeFindCostCenters(t *testing.T) {
	testCases := []struct {
		id   int64
		want string
	}{
		{25, fmt.Sprintf(string(employeeCostCentersEndpoint), 25)},
		{28, fmt.Sprintf(string(employeeCostCentersEndpoint), 28)},
	}

	for _, tc := range testCases {
		testPath(t, tc.want, func(c *Client) error {
			_, err := c.Employee.FindCostCenters(context.Background(), tc.id)
			return err
		})
	}

	testMethod(t, http.MethodGet, func(c *Client) error {
		_, err := c.Employee.FindCostCenters(context.Background(), 20)
		return err
	})

	testResponseBody(t, [][]byte{
		[]byte(`[{"id":77045,"name":"IT"}]`),
	}, func(c *Client) (interface{}, error) {
		return c.Employee.FindCostCenters(context.Background(), 10)
	})
}

func TestEmployeeFindCostCenterError(t *testing.T) {
	testError(t, func(c *Client) error {
		_, err := c.Employee.FindCostCenters(context.Background(), 10)
		return err
	})
}

func TestEmployeeUpdateCostCenters(t *testing.T) {
	testCases := []struct {
		id   int64
		want string
	}{
		{25, fmt.Sprintf(string(employeeCostCentersEndpoint), 25)},
		{28, fmt.Sprintf(string(employeeCostCentersEndpoint), 28)},
	}

	for _, tc := range testCases {
		testPath(t, tc.want, func(c *Client) error {
			_, err := c.Employee.UpdateCostCenters(context.Background(), tc.id, []int64{})
			return err
		})
	}

	testMethod(t, http.MethodPatch, func(c *Client) error {
		_, err := c.Employee.UpdateCostCenters(context.Background(), 20, []int64{})
		return err
	})

	testResponseBody(t, [][]byte{
		[]byte(`[100,200]`),
		[]byte(`[100]`),
	}, func(c *Client) (interface{}, error) {
		return c.Employee.UpdateCostCenters(context.Background(), 20, []int64{})
	})

	testRequestBody(t, []func(*Client) ([]byte, error){
		func(c *Client) (want []byte, err error) {
			want = []byte(`{"costCenterIDs":[200]}`)
			_, err = c.Employee.UpdateCostCenters(context.Background(), 20, []int64{200})
			return
		},
		func(c *Client) (want []byte, err error) {
			want = []byte(`{"costCenterIDs":[100,200]}`)
			_, err = c.Employee.UpdateCostCenters(context.Background(), 20, []int64{100, 200})
			return
		},
	})

}

func TestEmployeeUpdateCostCentersError(t *testing.T) {
	testError(t, func(c *Client) error {
		_, err := c.Employee.UpdateCostCenters(context.Background(), 10, nil)
		return err
	})
}
