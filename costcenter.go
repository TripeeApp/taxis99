package taxis99

import (
	"context"
	"net/http"
)

const (
	costCenterEndpoint endpoint = `costcenters`
)

// Hashset for allwed query params.
var ccFields = map[string]struct{}{
	"search": struct{}{},
	"limit":  struct{}{},
	"page":   struct{}{},
}

type CostCenter struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Enabled bool     `json:"enabled,omitempty"`
	Company *Company `json:"company,omitempty"`
}

type CostCenterService service

func (c *CostCenterService) Find(ctx context.Context, f Filter) ([]*CostCenter, error) {
	var costCenters []*CostCenter

	v := f.values(ccFields)

	err := c.client.Request(ctx, http.MethodGet, string(costCenterEndpoint.Query(v)), nil, &costCenters)
	if err != nil {
		return nil, err
	}

	return costCenters, nil
}

func (c *CostCenterService) Create(ctx context.Context, newCC CostCenter) (*CostCenter, error) {
	var cc *CostCenter

	err := c.client.Request(context.Background(), http.MethodPost, string(costCenterEndpoint), newCC, &cc)
	if err != nil {
		return nil, err
	}

	return cc, nil
}
