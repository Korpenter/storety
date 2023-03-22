package validators

import (
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateAuthorization(t *testing.T) {
	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				user: &models.User{
					Login:    "test",
					Password: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "empty login",
			args: args{
				user: &models.User{
					Login:    "",
					Password: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "empty password",
			args: args{
				user: &models.User{
					Login:    "test",
					Password: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthorization(tt.args.user)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
