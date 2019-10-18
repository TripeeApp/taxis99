package taxis99

import (
	"fmt"
	"net/url"
)

type endpoint string

func (e endpoint) Query(v url.Values) endpoint {
	if len(v) > 0 {
		e = endpoint(fmt.Sprintf("%s?%s", e, v.Encode()))
	}
	return e
}
