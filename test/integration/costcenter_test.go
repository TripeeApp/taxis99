package integration

import (
	"context"
	"testing"

	"github.com/mobilitee-smartmob/taxis99"
)

func TestCostCenter(t *testing.T) {
	newCC := taxis99.CostCenter{Name: randString(10, letterBytes)}

	cc, err := tx99.CostCenter.Create(context.Background(), newCC)
	if err != nil {
		t.Fatalf("Got error calling CostCenter.Create(%+v): '%s'; want nil.", newCC, err.Error())
	}
	// Removes the Cost Center once tests are done.
	defer func() {
		err := tx99.CostCenter.Remove(context.Background(), cc.ID)
		if err != nil {
			t.Errorf("Got error removing cost center %s: %s; want nil.", cc.Name, err.Error())
		}
	}()

	ccFilter := taxis99.Filter{"search": cc.Name}
	foundCC, err := tx99.CostCenter.Find(context.Background(), ccFilter)
	if err != nil {
		t.Errorf("Got error while calling CostCenter.Find(%+v): '%s'; want nil.", ccFilter, err.Error())
	}

	if len(foundCC) == 0 {
		t.Errorf("Got cost center %+v not found; want it to be found.", cc)
	}
}
