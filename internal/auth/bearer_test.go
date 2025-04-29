package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}

	_, err := GetBearerToken(header)

	if err == nil {
		t.Errorf("finding key in blank header, %v, error", err)
		return
	}

	header.Add("Authorization", "BigBadWolf")

	_, err = GetBearerToken(header)

	if err == nil {
		t.Errorf("finding key in invalid token, %v, error", err)
		return
	}

	header.Add("Authorization", "Bearer BigBadWolf")

	_, err = GetBearerToken(header)

	if err != nil {
		t.Errorf("did not find key in valid header, %v, error", err)
		return
	}
}
