package taxis99

import (
	"net/http"
)

// Transport is the RoundTripper to inject
// the authorization header for every request to the 99taxis API
type Transport struct {
	// ApiKey to be injected in the outgoing request
	APIKey string
	// CompanyID  
	CompanyID string
	// Base RoundTripper
	Base http.RoundTripper
}

// RoundTrip injects the authorization header with the token
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Do not modify the origin request
	req := cloneReq(r)
	// Inject the authorization headers
	req.Header.Set("x-api-key", t.APIKey)
	if t.CompanyID != "" {
		req.Header.Set("x-company-id", t.CompanyID)
	}

	return t.base().RoundTrip(req)
}


// base returns nil if base RoundTripper is nil
func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}

	return http.DefaultTransport
}

// cloneReq returns a shalow copy of the r *http.Request
// and deep copy of the Header map
func cloneReq(r *http.Request) *http.Request {
	// shalow copy the http.Request struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy the request headers struct
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}