package taxis99

import (
	"context"
	"net/http"
)

type Company struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

const companiesEndpoint endpoint = `companies`

type CompanyService service

func (c *CompanyService) Find(ctx context.Context) ([]*Company, error) {
	var companies []*Company
	err := c.client.Request(ctx, http.MethodGet, string(companiesEndpoint), nil, &companies)
	if err != nil {
		return nil, err
	}
	return companies, nil
}
