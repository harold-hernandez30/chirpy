package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")

	if bearer == "" {
		return "", fmt.Errorf("authorization not found")
	}

	fmt.Printf("LOG: bearer: %s\n", bearer)

	splitBearer := strings.Split(bearer, " ")

	fmt.Printf("LOG: splitBearer: %v\n", splitBearer)
	fmt.Printf("LOG: splitBearer[0]: %v\n", splitBearer[0])
	fmt.Printf("LOG: splitBearer[1]: %v\n", splitBearer[1])

	if len(splitBearer) == 2 {
		return splitBearer[1], nil
	}

	return "", fmt.Errorf("something went wrong. Invalid Authorization Bearer")
}