package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"github.com/halilylm/service/business/sys/auth"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	keyID := "random"
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	a, err := auth.New(keyID, &keyStore{pk: privateKey})
	if err != nil {
		t.Fatal(err)
	}
	t.Run("handling a single user", func(t *testing.T) {
		claims := auth.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "test project",
				Subject:   "random-user-id",
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				IssuedAt:  time.Now().UTC().Unix(),
			},
			Roles: []string{auth.RoleAdmin},
		}
		token, err := a.GenerateToken(claims)
		if err != nil {
			t.Fatal(err)
		}

		parsedClaims, err := a.ValidateToken(token)
		if err != nil {
			t.Fatal(err)
		}
		exp := len(claims.Roles)
		got := len(parsedClaims.Roles)
		assert.Equal(t, exp, got)
		assert.Equal(t, claims.Roles[0], parsedClaims.Roles[0])
	})
}

type keyStore struct {
	pk *rsa.PrivateKey
}

func (ks *keyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	return ks.pk, nil
}

func (ks *keyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	return &ks.pk.PublicKey, nil
}
