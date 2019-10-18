package taxis99

import (
	"net/url"
	"testing"
)

func TestEndpointQuery(t *testing.T) {
	testCases := []struct {
		endpoint endpoint
		values   url.Values
		want     endpoint
	}{
		{"employees", url.Values{"search": []string{"test"}}, "employees?search=test"},
		{"employees", url.Values{"costCenter": []string{"1"}}, "employees?costCenter=1"},
		{"employees", url.Values{}, "employees"},
	}

	for _, tc := range testCases {
		got := tc.endpoint.Query(tc.values)

		if got != tc.want {
			t.Errorf("Got endpoint '%s'; want '%s'.", got, tc.want)
		}
	}
}
