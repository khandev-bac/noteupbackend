package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TokenGeneration(userId uuid.UUID) (*Token, error) {
	accessclaims := TokenOUT{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "access",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshclaims := TokenOUT{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "refresh",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, accessclaims)
	accessToken, err := accessTokenClaim.SignedString([]byte(ACCESSTOKEN_KEY))
	if err != nil {
		return nil, err
	}
	refreshTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims)
	refreshToken, err := refreshTokenClaim.SignedString([]byte(REFRESHTOKEN_KEY))
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func VerifyAccessToken(token string) (*Payload, error) {
	decodedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid method")
		}
		return []byte(ACCESSTOKEN_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	claims := decodedToken.Claims.(jwt.MapClaims)
	if !decodedToken.Valid {
		return nil, errors.New("Invalid token,verification failed")
	}
	userIdstr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id claim missing or invalid")
	}
	userId, _ := uuid.Parse(userIdstr)
	return &Payload{
		UserId: userId,
	}, nil
}

func VerifyRefreshToken(token string) (*Payload, error) {
	decodedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid method")
		}
		return []byte(REFRESHTOKEN_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	claims := decodedToken.Claims.(jwt.MapClaims)
	if !decodedToken.Valid {
		return nil, errors.New("Invalid token,verification failed")
	}
	userIdstr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id claim missing or invalid")
	}
	userId, _ := uuid.Parse(userIdstr)
	return &Payload{
		UserId: userId,
	}, nil
}
