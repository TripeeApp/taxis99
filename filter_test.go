package taxis99

import (
	"net/url"
	"reflect"
	"testing"
)

func TestFiltersAdd(t *testing.T) {
	got := Filters{}
	got.Set("search", "test")

	if want := (Filters{"search": "test"}); !reflect.DeepEqual(got, want) {
		t.Errorf("Got Filters: %+v; want %+v.", got, want)
	}
}

func TestFiltersDel(t *testing.T) {
	got := Filters{"search": "test"}
	got.Del("search")

	if want := (Filters{}); !reflect.DeepEqual(got, want) {
		t.Errorf("Got Filters: %+v; want %+v.", got, want)
	}
}

func TestFilterValues(t *testing.T) {
	fields := map[string]struct{}{
		"search": struct{}{},
		"limit":  struct{}{},
	}
	f := Filters{
		"search":       "test",
		"costCenterId": "1",
	}
	got := f.values(fields)

	if want := (url.Values{"search": []string{"test"}}); !reflect.DeepEqual(got, want) {
		t.Errorf("Got url.Values: %+v; want %+v.", got, want)
	}
}
