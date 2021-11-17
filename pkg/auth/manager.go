package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager interface {
	NewJWT(userID string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) *Manager {
	return &Manager{signingKey: signingKey}
}

func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    userId,
		ExpiresAt: time.Now().Add(100000000000 * 60).Unix(),
	})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(tokenString string) (bool, interface{}) {
	if tokenString == "" {
		return false, ""
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	fmt.Printf("token.Claims: %v\n", token.Claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["iss"], claims["exp"])
		return true, claims["iss"]
	} else {
		fmt.Println(err)
		return false, ""
	}
}
