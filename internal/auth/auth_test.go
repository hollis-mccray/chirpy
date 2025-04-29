package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHash(t *testing.T) {
	password := "DreadLordMordred"
	hashCode, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("Unable to generate hashcode, %v, error", err)
		return
	}

	err = CheckPasswordHash(string(hashCode), password)
	if err != nil {
		t.Errorf("Hashcode doesnot match with password, %v, %v, %v, error", hashCode, password, err)
	}
}

func TestMakeJWT(t *testing.T) {
	tokenSecret := "DreadLordMordred"
	after, err := time.ParseDuration("5s")
	if err != nil {
		t.Errorf("error parsing duration, %v, error", err)
		return
	}
	_, err = MakeJWT(
		uuid.New(),
		tokenSecret,
		after,
	)

	if err != nil {
		t.Errorf("error creating jwt, %v, error", err)
		return
	}
}

func TestValidateJWT(t *testing.T) {
	tokenSecret := "DreadLordMordred"
	after, err := time.ParseDuration("5h")
	if err != nil {
		t.Errorf("error parsing duration, %v, error", err)
		return
	}
	oldId := uuid.New()
	token, err := MakeJWT(
		oldId,
		tokenSecret,
		after,
	)

	if err != nil {
		t.Errorf("error creating jwt, %v, error", err)
		return
	}

	newId, err := ValidateJWT(token, tokenSecret)

	if err != nil {
		t.Errorf("error validating jwt, %v, error", err)
		return
	}

	if newId != oldId {
		t.Errorf("uuids not equal, error")
	}
}

func TestInvalidSecret(t *testing.T) {
	validSecret := "DreadLordMordred"
	invalidSecret := "BadKitty"
	after, err := time.ParseDuration("5h")
	if err != nil {
		t.Errorf("error parsing duration, %v, error", err)
		return
	}
	oldId := uuid.New()
	token, err := MakeJWT(
		oldId,
		validSecret,
		after,
	)

	if err != nil {
		t.Errorf("error creating jwt, %v, error", err)
		return
	}

	_, err = ValidateJWT(token, invalidSecret)

	if err == nil {
		t.Errorf("failure to invalidate invalid secret, %v, error", err)
		return
	}
}

func TestExpiredToken(t *testing.T) {
	tokenSecret := "DreadLordMordred"
	after, err := time.ParseDuration("1s")
	if err != nil {
		t.Errorf("error parsing duration, %v, error", err)
		return
	}
	oldId := uuid.New()
	token, err := MakeJWT(
		oldId,
		tokenSecret,
		after,
	)

	if err != nil {
		t.Errorf("error creating jwt, %v, error", err)
		return
	}

	time.Sleep(5 * time.Second)

	_, err = ValidateJWT(token, tokenSecret)

	if err == nil {
		t.Errorf("failure to invalidate expired token, %v, error", err)
		return
	}
}
