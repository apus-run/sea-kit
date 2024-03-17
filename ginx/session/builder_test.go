package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/apus-run/sea-kit/ginx"
)

func TestBuilder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := NewMockProvider(ctrl)
	p.EXPECT().NewSession(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx *ginx.Context, uid int64, jwtData map[string]string,
			sessData map[string]any) (Session, error) {
			return &MemorySession{data: sessData,
				claims: Claims{Uid: uid, Data: jwtData}}, nil
		})
	sess, err := NewSessionBuilder(new(ginx.Context), 123).
		SetProvider(p).
		SetJwtData(map[string]string{"jwt": "true"}).
		SetSessData(map[string]any{"session": "true"}).
		Build()
	require.NoError(t, err)
	assert.Equal(t, &MemorySession{
		data: map[string]any{"session": "true"},
		claims: Claims{
			Uid:  123,
			Data: map[string]string{"jwt": "true"},
		},
	}, sess)
}
