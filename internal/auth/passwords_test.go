package auth

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	plaintext := "ldkasjfoiejrfolwdfj@#$@#"
	hashedPassword, err := HashPassword(plaintext)

	fmt.Println("Test running for TestHashPassword")
	fmt.Printf("plaintext: %s\n", plaintext)
	fmt.Printf("hashedPassword: %s\n", hashedPassword)
	if err != nil {
		t.Errorf("could not hash password: %v", err)
	}

	checkPasswordHashErr := CheckPasswordHash(plaintext, hashedPassword)

	if checkPasswordHashErr != nil {
		t.Errorf("check password hash failed: %v\n", checkPasswordHashErr)
	}

	fmt.Println("Password matches!")
}