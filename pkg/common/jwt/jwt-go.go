package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/erviangelar/go-users-api/pkg/common/config"
	"github.com/erviangelar/go-users-api/pkg/common/models"
	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	KeyType  string
	jwt.StandardClaims
}
type JWTRefreshClaim struct {
	ID string `json:"id"`
	// CustomKey string
	KeyType string
	jwt.StandardClaims
}

// generate refresh token
func GenerateRefreshToken(user *models.User) (string, error) {
	configs := config.NewConfigurations()
	userID := strconv.Itoa((int)(user.ID))
	// cusKey := GenerateCustomKey(userID, configs.TokenHash)
	tokenType := "refresh"

	claims := JWTRefreshClaim{
		userID,
		// cusKey,
		tokenType,
		jwt.StandardClaims{
			Issuer: "auth.service",
		},
	}

	signBytes, err := os.ReadFile(configs.RefreshTokenPrivateKeyPath)
	if err != nil {
		return "", errors.New("could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", errors.New("could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(signKey)
}

// generate access token
func GenerateAccessToken(user *models.User) (string, error) {

	configs := config.NewConfigurations()
	userID := strconv.Itoa((int)(user.ID))
	tokenType := "access"

	claims := JWTClaim{
		userID,
		user.Username,
		user.Name,
		user.Role,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(configs.JwtExpiration)).Unix(),
			Issuer:    "auth.service",
		},
	}

	signBytes, err := os.ReadFile(configs.AccessTokenPrivateKeyPath)
	if err != nil {
		return "", errors.New("could not generate access token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", errors.New("could not generate access token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// generate custom key
func GenerateCustomKey(userID string, tokenHash string) string {

	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userID))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// validate access token
func ValidateAccessToken(tokenString string) (string, error) {

	configs := config.NewConfigurations()
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method in auth token")
		}
		verifyBytes, err := os.ReadFile(configs.AccessTokenPublicKeyPath)
		if err != nil {
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid || claims.ID == "" || claims.KeyType != "access" {
		return "", errors.New("invalid token: authentication failed")
	}
	return claims.Role, nil
}

// validate refresh token
func ValidateRefreshToken(tokenString string) (string, error) {

	configs := config.NewConfigurations()
	token, err := jwt.ParseWithClaims(tokenString, &JWTRefreshClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method in auth token")
		}
		verifyBytes, err := os.ReadFile(configs.RefreshTokenPublicKeyPath)
		if err != nil {
			return nil, errors.New("invalid token: authentication failed")
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			return nil, errors.New("invalid token: authentication failed")
		}

		return verifyKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*JWTRefreshClaim)
	fmt.Println(claims.ID)
	fmt.Println(claims.KeyType)
	if !ok || !token.Valid || claims.ID == "" || claims.KeyType != "refresh" {
		return "", errors.New("invalid token: authentication failed")
	}
	return claims.ID, nil
}
