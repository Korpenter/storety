package token

import (
	"fmt"
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

// JWTAuth implements JWT authentication flow.
type JWTAuth struct {
	cfg *config.Config
}

// NewJwtAuth configures and returns a JWT authentication instance.
func NewJwtAuth(i *do.Injector) *JWTAuth {
	cfg := do.MustInvoke[*config.Config](i)
	return &JWTAuth{
		cfg: cfg,
	}
}

// Verify checks the validity of the given JWT token and returns the UUID associated with it.
// Returns an error if the token is invalid or expired.
func (a *JWTAuth) Verify(token string) (uuid.UUID, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.JWTAuthKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}
	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// GenerateTokenPair returns both an access token and a refresh token for the given UUID.
func (a *JWTAuth) GenerateTokenPair(id, sessionID uuid.UUID) (string, string, error) {
	access, err := a.createJWT(id)
	if err != nil {
		return "", "", err
	}
	refresh, err := a.createRefreshJWT(sessionID)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// createJWT returns an access token for the provided UUID.
func (a *JWTAuth) createJWT(id uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().
			Add(time.Duration(a.cfg.JWTAuthLifeTimeHours) * time.Hour)),
		ID: id.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.JWTAuthKey))
	return tokenString, err
}

// createRefreshJWT returns a refresh token for the provided UUID.
func (a *JWTAuth) createRefreshJWT(id uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().
			Add(time.Duration(a.cfg.JWTRefreshLifeTimeHours) * time.Hour)),
		ID: id.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.JWTAuthKey))
	return tokenString, err
}
