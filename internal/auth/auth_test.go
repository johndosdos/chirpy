package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

const TOKEN_SECRET = "secret-secret-too"

var tokenString string
var err error

func TestMakeJWT(t *testing.T) {
	tokenString, err = MakeJWT(uuid.New(), TOKEN_SECRET, 2*time.Second)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	fmt.Println("Token string: ", tokenString)
}

func TestValidateJWT(t *testing.T) {
	tokenString, err = MakeJWT(uuid.New(), TOKEN_SECRET, 2*time.Second)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	fmt.Println("\nToken string: ", tokenString)

	userID, err := ValidateJWT(tokenString, TOKEN_SECRET)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	fmt.Println(userID)
}
