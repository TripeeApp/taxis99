package taxis99

import (
	"net/url"
)

// Filter is used to filter the requests to the API.
type Filter map[string]string

// Set updates the key value with the val.
// If there's already a value, it will be replaced.
func (f Filter) Set(key, val string) {
	f[key] = val
}

func (f Filter) Del(key string) {
	delete(f, key)
}

// Values returns a url.Values mapped between Filter and values. Used internally only.
func (f Filter) values(fields map[string]struct{}) url.Values {
	vals := url.Values{}
	for k, v := range f {
		if _, ok := fields[k]; ok {
			vals.Add(k, v)
		}
	}
	return vals
}
