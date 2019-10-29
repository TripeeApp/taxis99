package taxis99

import (
	"net/http"
)

const (
	headerAPIKey    = "X-Api-Key"
	headerCompanyID = "X-Company-Id"
)

type companyIDKey struct{}

var CompanyID companyIDKey

// Transport is the RountTripper for injecting
// the authorization header for every request to
// 99 taxis API.
type Transport struct {
	// Key is the string to be injected to the
	// outgoing request header
	Key string

	// CompanyID is the string related to a company
	// that will be injected to the outgoing
	// requests if it's not empty.
	CompanyID string

	// Base is the base RoundTripper to make HTTP request.
	Base http.RoundTripper
}

// RoundTrip injects the Authorization Header with the key
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	// We should not modify the origin request
	// per RoundTripper contract. See
	// https://golang.org/pkg/net/http/#RoundTripper
	req := cloneReq(r)
	// Injects the Authorization Header
	req.Header.Set(headerAPIKey, t.Key)

	if t.CompanyID != "" {
		req.Header.Set(headerCompanyID, t.CompanyID)
	}

	if cid, ok := req.Context().Value(CompanyID).(string); ok {
		req.Header.Set(headerCompanyID, cid)
	}

	return t.base().RoundTrip(req)
}

// base returns the base RoundTripper or the http default transport.
func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}

	return http.DefaultTransport
}

// cloneReq returns a clone of the *http.Request.
// the clone is a shallow copy of the struct and its Header map.
func cloneReq(r *http.Request) *http.Request {
	// shalow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
