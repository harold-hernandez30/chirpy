package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	
	bearer := headers.Get("Authorization")

	if bearer == "" {
		return "", fmt.Errorf("authorization not found")
	}

	fmt.Printf("LOG: bearer: %s\n", bearer)

	splitBearer := strings.Split(bearer, " ")

	if splitBearer[0] != "ApiKey" {
		return "", fmt.Errorf("expecting ApiKey Authorization")
	}

	if len(splitBearer[1]) == 0 {
		return "", fmt.Errorf("invalid ApiKey")
	}


	if len(splitBearer) == 2 {
		return splitBearer[1], nil
	}

	return "", fmt.Errorf("something went wrong. Invalid Authorization ApiKey")

}