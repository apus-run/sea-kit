package auth

import (
	"context"
	"testing"

	"github.com/apus-run/sea-kit/redisx"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	*redisx.Helper
}

func TestAuthenticator(t *testing.T) {
	tests := []struct {
		name     string
		app      string
		token    string
		strict   bool
		hasError bool
	}{
		{
			name:     "strict=false",
			strict:   false,
			hasError: false,
		},
		{
			name:     "strict=true",
			strict:   true,
			hasError: true,
		},
		{
			name:     "strict=true,with token",
			app:      "foo",
			token:    "bar",
			strict:   true,
			hasError: false,
		},
		{
			name:     "strict=true,with error token",
			app:      "foo",
			token:    "error",
			strict:   true,
			hasError: true,
		},
	}

	ctx := context.Background()
	h := &Client{redisx.NewHelper()}
	client, err := h.GetClient(redisx.WithRedisConfig(func(options *redisx.RedisConfig) {
		options.Addr = "localhost:16379"
		options.DB = 0
		options.Username = "root"
	}))
	if err != nil {
		t.Fatal(err)
	}

	// 检测数据库是否可以连接
	cmd := client.Ping(ctx)
	if cmd.Err() != nil {
		t.Fatal(cmd.Err())
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.app) > 0 {
				assert.Nil(t, client.HSet(context.Background(), "apps", test.app, test.token).Err())
				defer client.HDel(context.Background(), "apps", test.app)
			}

			authenticator, err := NewAuthenticator(client, "apps", test.strict)
			assert.Nil(t, err)
			assert.NotNil(t, authenticator.Authenticate(context.Background()))
			md := metadata.New(map[string]string{})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			assert.NotNil(t, authenticator.Authenticate(ctx))
			md = metadata.New(map[string]string{
				"app":   "",
				"token": "",
			})
			ctx = metadata.NewIncomingContext(context.Background(), md)
			assert.NotNil(t, authenticator.Authenticate(ctx))
			md = metadata.New(map[string]string{
				"app":   "foo",
				"token": "bar",
			})
			ctx = metadata.NewIncomingContext(context.Background(), md)
			err = authenticator.Authenticate(ctx)
			if test.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
