package infrastructure

import (
	"errors"
	"fmt"
	"time"
	"trainer/internal/application"
	"trainer/internal/domain/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JwtManager struct {
	secretKey string
	ttl       time.Duration
}

func (m *JwtManager) Generate(claim application.TokenClaim) (string, error) {
	claims := jwt.MapClaims{
		"user_id": claim.UserID.String(),
		"role":    string(claim.Role),
		"exp":     time.Now().Add(m.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JwtManager) Parse(tokenStr string) (*application.TokenClaim, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, errors.New("token expired")
		}
	} else {
		return nil, errors.New("exp claim missing or invalid")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id claim is not a string")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("role claim is not a string")
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	return &application.TokenClaim{
		UserID: parsedUUID,
		Role:   user.Role(role),
	}, nil
}
