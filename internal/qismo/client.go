package qismo

import (
	"context"
	"io"
)

type httpClient interface {
	Call(ctx context.Context, method, url string, body io.Reader, headers map[string]string, response interface{}) error
}
