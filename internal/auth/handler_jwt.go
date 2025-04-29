package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: &jwt.NumericDate{
			Time: time.Now().UTC(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().UTC().Add(expiresIn),
		},
		Subject: userID.String(),
	}).SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to parse token : %v", err)
	}

	idString, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to extract id from claim : %v", err)
	}

	return uuid.Parse(idString)
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Values("Authorization")
	if len(auth) == 0 {
		return "", fmt.Errorf("authorization key not found")
	}

	for _, value := range auth {
		token, found := strings.CutPrefix(value, "Bearer ")
		if found {
			return token, nil
		}
	}

	return "", fmt.Errorf("bearer token not found")
}
