package helpers

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestCheckMDValue(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		key    string
		want   string
		wantOk bool
	}{

		{
			name: "Valid test",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.New(map[string]string{"refresh_token": "expiredToken"})),
			key:    "refresh_token",
			want:   "expiredToken",
			wantOk: true,
		},
		{
			name:   "Invalid test",
			ctx:    context.Background(),
			key:    "refresh_token",
			want:   "",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md, ok := CheckMDValue(tt.ctx, tt.key)
			if tt.wantOk {
				require.True(t, ok)
				require.Equal(t, tt.want, md)
			}
		})
	}
}
