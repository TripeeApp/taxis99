package taxis99

import (
	"net/url"
	"reflect"
	"testing"
)

func TestFilterAdd(t *testing.T) {
	got := Filter{}
	got.Set("search", "test")

	if want := (Filter{"search": "test"}); !reflect.DeepEqual(got, want) {
		t.Errorf("Got Filter: %+v; want %+v.", got, want)
	}
}

func TestFilterDel(t *testing.T) {
	got := Filter{"search": "test"}
	got.Del("search")

	if want := (Filter{}); !reflect.DeepEqual(got, want) {
		t.Errorf("Got Filter: %+v; want %+v.", got, want)
	}
}

func TestFilterValues(t *testing.T) {
	fields := map[string]struct{}{
		"search": struct{}{},
		"limit":  struct{}{},
	}
	f := Filter{
		"search":       "test",
		"costCenterId": "1",
	}
	got := f.values(fields)

	if want := (url.Values{"search": []string{"test"}}); !reflect.DeepEqual(got, want) {
		t.Errorf("Got url.Values: %+v; want %+v.", got, want)
	}
}
