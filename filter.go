package taxis99

import (
	"net/url"
)

// Filter is used to filter the requests to the API.
type Filters map[string]string

// Set updates the key value with the val.
// If there's already a value, it will be replaced.
func (f Filters) Set(key, val string) {
	f[key] = val
}

func (f Filters) Del(key string) {
	delete(f, key)
}

// Values returns a url.Values mapped between Filter and values. Used internally only.
func (f Filters) values(fields map[string]struct{}) url.Values {
	vals := url.Values{}
	for k, filter := range f {
		if _, ok := fields[k]; ok {
			vals.Add(k, filter)
		}
	}
	return vals
}
