package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	refUUID, uuidErr := uuid.NewRandom()

	if uuidErr != nil {
		t.Error("unable to create UUID")
	}

	secret := "whysosecret"
	expiresIn := 3 * time.Second
	tokenString, makeJwtErr := MakeJWT(refUUID, secret, expiresIn)

	if makeJwtErr != nil {
		t.Errorf("could not make jwt: %s\n", makeJwtErr)
		return
	}

	uuidFromClaim, validateJWTErr := ValidateJWT(tokenString, secret)

	if validateJWTErr != nil {
		t.Errorf("validation failed: %s\n", validateJWTErr)
		return
	}


	fmt.Println("Test running for TestMakeJWT")
	fmt.Printf("secret: %s\n", secret)
	fmt.Printf("tokenString: %s\n", tokenString)
	fmt.Printf("uuid reference: %s\n", refUUID.String())
	fmt.Printf("uuid from claim: %s\n", uuidFromClaim.String())
	
	if uuidFromClaim.String() != refUUID.String() {
		t.Error("UUID from validateJWT not equal\n")
		return
	}

	fmt.Println("PASS!")
}