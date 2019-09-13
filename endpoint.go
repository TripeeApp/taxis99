package taxis99

import (
	"context"
	"fmt"
)

type endpoint string

// String format the endpoint as a string
func (e endpoint) String(ctx context.Context) string {
	return fmt.Sprintf("v2/%s", string(e))
}