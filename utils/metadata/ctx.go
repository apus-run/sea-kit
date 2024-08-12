package metadata

import (
	"context"
	"net/http"
	"strings"
)

func FromRequest(r *http.Request) context.Context {
	ctx := context.Background()
	md := make(Metadata)
	for k, v := range r.Header {
		md[k] = strings.Join(v, ",")
	}
	return NewContext(ctx, md)
}
