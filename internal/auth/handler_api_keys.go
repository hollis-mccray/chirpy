package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Values("Authorization")
	if len(auth) == 0 {
		return "", fmt.Errorf("api key not found")
	}

	for _, value := range auth {
		token, found := strings.CutPrefix(value, "ApiKey ")
		if found {
			return token, nil
		}
	}

	return "", fmt.Errorf("API key not found")

}
