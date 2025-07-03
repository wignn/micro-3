package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/wignn/micro-3/auth/model"
)

var (
	AccessSecretKey  []byte
	RefreshSecretKey []byte
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
)

type Config struct {
	ACCESS_SECRET_KEY  string `envconfig:"ACCESS_SECRET_KEY"`
	REFRESH_SECRET_KEY string `envconfig:"REFRESH_SECRET_KEY"`
}

const (
	DefaultAccessTokenTTL  = 15 * time.Minute
	DefaultRefreshTokenTTL = 30 * 24 * time.Hour
)

var cfg Config


func InitJWTConfig() {
	if err := envconfig.Process("", &cfg); err != nil {
		fmt.Println("Failed to process environment variables:", err)
	}

	AccessSecretKey = []byte(cfg.ACCESS_SECRET_KEY)
	RefreshSecretKey = []byte(cfg.REFRESH_SECRET_KEY)
	AccessTokenTTL = DefaultAccessTokenTTL
	RefreshTokenTTL = DefaultRefreshTokenTTL

	log.Printf("RefreshSecretKey=%s", string(RefreshSecretKey))
	log.Printf("AccessSecretKey=%s", string(AccessSecretKey))
	log.Printf("JWT Config initialized: AccessTTL=%s, RefreshTTL=%s", AccessTokenTTL, RefreshTokenTTL)
	if AccessSecretKey == nil || RefreshSecretKey == nil {
		log.Fatal("JWT secret keys must be set")
	}
}

// GenerateToken creates a signed access and refresh token
func GenerateToken(email string) (*model.Token, error) {
	now := time.Now()
	accessExpiry := now.Add(AccessTokenTTL)
	refreshExpiry := now.Add(RefreshTokenTTL)

	accessClaims := jwt.MapClaims{
		"email": email,
		"exp":   accessExpiry.Unix(),
		"iat":   now.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(AccessSecretKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"email": email,
		"exp":   refreshExpiry.Unix(),
		"iat":   now.Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(RefreshSecretKey)
	if err != nil {
		return nil, err
	}

	return &model.Token{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		ExpiresAt:    uint64(accessExpiry.Unix()),
	}, nil
}
func ValidateRefreshToken(tokenStr string) (string, error) {
	log.Printf("Validating refresh token: %s", tokenStr)
	log.Printf("RefreshSecretKey=%s", string(cfg.REFRESH_SECRET_KEY))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		log.Printf("Token algorithm: %v", token.Header["alg"])
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// FIX: convert secret key to []byte
		return []byte(cfg.REFRESH_SECRET_KEY), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return "", err
	}
	if !token.Valid {
		log.Println("Token is not valid")
		return "", errors.New("token not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to convert token claims to MapClaims")
		return "", errors.New("could not parse claims")
	}

	log.Printf("Parsed claims: %+v", claims)

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("Email claim not found or not a string")
		return "", errors.New("email claim not found")
	}

	log.Printf("Extracted email from token: %s", email)
	return email, nil
}

// ParseAccessToken parses and validates an access token
func ParseAccessToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing algorithm")
		}
		return AccessSecretKey, nil
	})

	log.Printf("Parsing access token: %s", tokenStr)
	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", errors.New("invalid email in token")
	}

	if exp, ok := claims["exp"].(float64); !ok || time.Now().Unix() > int64(exp) {
		return "", errors.New("access token expired")
	}

	return email, nil
}
