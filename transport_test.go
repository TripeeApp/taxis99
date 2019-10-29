package taxis99

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRoundTripper func(r *http.Request) (*http.Response, error)

func (rt testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

func TestTransportRoundTrip(t *testing.T) {
	testCases := []struct {
		wantAPIKey    string
		wantCompanyID string
	}{
		{"x-abc-key", ""},
		{"x-abc-key", "abc"},
	}

	for _, tc := range testCases {
		var reqSent *http.Request

		rt := testRoundTripper(func(r *http.Request) (*http.Response, error) {
			reqSent = r
			return nil, nil
		})

		tr := &Transport{tc.wantAPIKey, tc.wantCompanyID, rt}

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		_, err := tr.RoundTrip(req)
		if err != nil {
			t.Fatalf("Got error calling Transport.RoundTrip: %s; want it to be nil.", err.Error())
		}

		if got := reqSent.Header.Get(headerAPIKey); got != tc.wantAPIKey {
			t.Errorf("Got %s Header: %s; want %s.", headerAPIKey, got, tc.wantAPIKey)
		}

		if got := reqSent.Header.Get(headerCompanyID); got != tc.wantCompanyID {
			t.Errorf("Got %s Header: %s; want %s.", headerCompanyID, got, tc.wantCompanyID)
		}
	}
}

func TestTransportRoundTripContext(t *testing.T) {
	testCases := []struct {
		context            context.Context
		transportCompanyID string
		wantCompanyID      string
	}{
		{context.WithValue(context.Background(), CompanyID, "123"), "abc", "123"},
		{context.WithValue(context.Background(), CompanyID, 1), "abc", "abc"},
		{context.Background(), "abc", "abc"},
	}

	for _, tc := range testCases {
		var reqSent *http.Request

		rt := testRoundTripper(func(r *http.Request) (*http.Response, error) {
			reqSent = r
			return nil, nil
		})

		tr := &Transport{"", tc.transportCompanyID, rt}

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		req = req.WithContext(tc.context)

		_, err := tr.RoundTrip(req)
		if err != nil {
			t.Fatalf("Got error calling Transport.RoundTrip: %s; want it to be nil.", err.Error())
		}

		if got := reqSent.Header.Get(headerCompanyID); got != tc.wantCompanyID {
			t.Errorf("Got %s Header: %s; want %s.", headerCompanyID, got, tc.wantCompanyID)
		}
	}
}

func TestTransportRoundTripError(t *testing.T) {
	rt := testRoundTripper(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("Error")
	})

	tr := &Transport{"", "", rt}

	_, err := tr.RoundTrip(httptest.NewRequest(http.MethodGet, "/", nil))
	if err == nil {
		t.Error("got error nil; want not nil")
	}
}
