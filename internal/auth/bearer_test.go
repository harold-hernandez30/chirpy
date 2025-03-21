package auth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	header_bearer := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIzMTQwOWNiYi0wNTc4LTRjYzgtYjRhNC1mNzk0ZjJhNWQ1NzUiLCJleHAiOjE3NDI1NjkxOTUsImlhdCI6MTc0MjU2OTE5NX0.FN99mI3R6ddMJrxA1Om94bWzgEmvYArHfQ4FjGyumuA"
	headers.Add("Authorization", header_bearer)
	bearerToken, getBearerTokenErr := GetBearerToken(headers)

	if getBearerTokenErr != nil {
		t.Errorf("something went wrong. Could not get bearer token: %s\n", bearerToken)
	}

	fmt.Println("===========================")
	fmt.Println("TestGetBearerToken")
	fmt.Println("===========================")
	fmt.Printf("Authorization (provided): %s\n", header_bearer)
	fmt.Printf("Authorization (extracted): %s\n", bearerToken)

	if bearerToken != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIzMTQwOWNiYi0wNTc4LTRjYzgtYjRhNC1mNzk0ZjJhNWQ1NzUiLCJleHAiOjE3NDI1NjkxOTUsImlhdCI6MTc0MjU2OTE5NX0.FN99mI3R6ddMJrxA1Om94bWzgEmvYArHfQ4FjGyumuA" {
		t.Errorf("bearerToken not equal")
		fmt.Println("FAILED!")
	} else {
		fmt.Println("PASSED!")
	}
}
