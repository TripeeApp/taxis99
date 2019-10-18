package taxis99

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestCostCenterFind(t *testing.T) {
	payload := []byte(`[{"id":"123","name":"TI","enabled":true,"company":{"id":"1234","name":"Mobilitee"}}]`)

	t.Run("ResponseBody", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Write(payload)
		}

		c, svr := newMockServer(nil, handler)
		defer svr.Close()

		got, err := c.CostCenter.Find(context.Background(), nil)
		if err != nil {
			t.Fatalf("Got error calling CostCenter.Find '%s'; want nil.", err.Error())
		}

		want := []*CostCenter{
			&CostCenter{
				ID:      "123",
				Name:    "TI",
				Enabled: true,
				Company: Company{
					ID:   "1234",
					Name: "Mobilitee",
				},
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got CostCenter %+v; want %+v.", got, want)
		}
	})

	t.Run("Path", func(t *testing.T) {
		var got string
		handler := func(w http.ResponseWriter, r *http.Request) {
			got = r.URL.Path
			w.Write(payload)
		}

		c, svr := newMockServer(nil, handler)
		defer svr.Close()

		c.CostCenter.Find(context.Background(), nil)

		if want := fmt.Sprintf("/%s", findCostCenter); got != want {
			t.Errorf("Got path '%s'; want '%s'.", got, want)
		}
	})

	t.Run("Method", func(t *testing.T) {
		var got string
		handler := func(w http.ResponseWriter, r *http.Request) {
			got = r.Method
			w.Write(payload)
		}

		c, svr := newMockServer(nil, handler)
		defer svr.Close()

		c.CostCenter.Find(context.Background(), nil)

		if want := http.MethodGet; got != want {
			t.Errorf("Got method '%s'; want '%s'.", got, want)
		}
	})

	t.Run("Query", func(t *testing.T) {
		testCases := []struct {
			filters Filters
		}{
			{Filters{"search": "123"}},
			{Filters{"search": "124", "limit": "100"}},
		}

		for _, tc := range testCases {
			var got url.Values
			handler := func(w http.ResponseWriter, r *http.Request) {
				got = r.URL.Query()
				w.Write(payload)
			}

			c, svr := newMockServer(nil, handler)
			defer svr.Close()

			c.CostCenter.Find(context.Background(), tc.filters)

			if want := tc.filters.values(ccFields); !reflect.DeepEqual(got, want) {
				t.Errorf("Got query values '%+v'; want '%+v'.", got, want)
			}
		}
	})
}

func TestCostCenterFindError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("1"))
	}

	c, svr := newMockServer(nil, handler)
	defer svr.Close()

	if _, err := c.CostCenter.Find(context.Background(), nil); err == nil {
		t.Error("Got error nil; want it not nil.")
	}

}
