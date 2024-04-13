package util

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	SessionId string
	Email     string
	Role      string
	jwt.RegisteredClaims
}

func GenerateToken(userId, userType, email, sessId string, exp time.Time) (tokenStr string, sessionId string, err error) {
	now := time.Now()
	claims := CustomClaims{
		SessionId: sessId,
		Email:     email,
		Role:      userType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://imaxinacion.com",
			Subject:   userId,
			Audience:  []string{"juggernaut-api"},
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte("MAXIMINA_TOKEN"))
	if err != nil {
		return "", "", fmt.Errorf("token.signedstring err: %w", err)
	}
	return signedString, sessId, nil

}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	var p jwt.Parser

	TokenClaims, _, err := p.ParseUnverified(tokenStr, &CustomClaims{})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New(" invalid signature error")
		}
		return nil, err
	}

	if claims, ok := TokenClaims.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	return nil, errors.New("error converting")
}

func ValidateToken(tokenStr string) (*CustomClaims, error) {
	Keyfunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("MAXIMINA_TOKEN"), nil
	}

	TokenClaims, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, Keyfunc)
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New(" invalid signature error")
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New(" token expired !")
		}
		return nil, err
	}

	if claims, ok := TokenClaims.Claims.(*CustomClaims); ok && TokenClaims.Valid {
		return claims, nil
	}
	return nil, errors.New("error converting")
}
