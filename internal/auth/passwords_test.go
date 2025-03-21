package auth

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {

	fmt.Println("===========================")
	fmt.Println("TestHashPassword")
	fmt.Println("===========================")
	plaintext := "ldkasjfoiejrfolwdfj@#$@#"
	hashedPassword, err := HashPassword(plaintext)

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