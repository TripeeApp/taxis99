package taxis99

import (
	"context"
	"net/http"
)

const (
	findCostCenter endpoint = `costcenters`
)

// Hashset for allwed query params.
var ccFields = map[string]struct{}{
	"search": struct{}{},
	"limit":  struct{}{},
	"page":   struct{}{},
}

type CostCenter struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Enabled bool    `json:"enabled"`
	Company Company `json:"company"`
}

type CostCenterService service

func (c *CostCenterService) Find(ctx context.Context, f Filters) ([]*CostCenter, error) {
	var costCenters []*CostCenter

	v := f.values(ccFields)

	err := c.client.Request(ctx, http.MethodGet, string(findCostCenter.Query(v)), nil, &costCenters)
	if err != nil {
		return nil, err
	}

	return costCenters, nil
}
