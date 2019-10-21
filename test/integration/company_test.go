package integration

import (
	"context"
	"testing"
)

func TestCompany(t *testing.T) {
	companies, err := tx99.Company.Find(context.Background())
	if err != nil {
		t.Errorf("Got error while calling Company.Find(): '%s'; want nil.", err.Error())
	}

	if len(companies) == 0 {
		t.Error("Got no company found; want it to be found.")
	}
}
