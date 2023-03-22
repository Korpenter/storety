package token

import (
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTokenAuth_GenerateTokenPair(t *testing.T) {
	cfg := &config.Config{
		JWTAuthKey:              "testKey",
		JWTAuthLifeTimeHours:    12,
		JWTRefreshLifeTimeHours: 240,
	}
	jwtAuth := JWTAuth{cfg: cfg}
	id, err := uuid.NewRandom()
	require.NoError(t, err)
	authToken, refreshToken, err := jwtAuth.GenerateTokenPair(id, id)
	require.NoError(t, err)
	require.NotEmpty(t, authToken)
	require.NotEmpty(t, refreshToken)
}

func TestTokenAuth_Verify(t *testing.T) {
	tests := []struct {
		name                    string
		JWTAuthLifeTimeHours    int
		JWTRefreshLifeTimeHours int
		wantErr                 bool
	}{
		{
			name:                    "Verify valid token",
			JWTAuthLifeTimeHours:    12,
			JWTRefreshLifeTimeHours: 240,
			wantErr:                 false,
		},
		{
			name:                    "Verify invalid token",
			JWTAuthLifeTimeHours:    0,
			JWTRefreshLifeTimeHours: 0,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				JWTAuthKey:              "testKey",
				JWTAuthLifeTimeHours:    tt.JWTAuthLifeTimeHours,
				JWTRefreshLifeTimeHours: tt.JWTRefreshLifeTimeHours,
			}
			auth := JWTAuth{cfg: cfg}
			id, err := uuid.NewRandom()
			require.NoError(t, err)
			tokenStr, err := auth.createJWT(id)
			require.NoError(t, err)
			require.NotEmpty(t, tokenStr)
			uid, err := auth.Verify(tokenStr)
			if tt.wantErr {
				require.Error(t, err)
				require.EqualValues(t, uuid.Nil, uid)
				return
			}
			require.NoError(t, err)
			require.Equal(t, id.String(), uid.String())
		})
	}
}
