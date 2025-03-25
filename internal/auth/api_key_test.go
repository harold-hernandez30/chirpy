package auth

import (
	"fmt"
	"net/http"
	"testing"
)



func TestGetApiKey(t *testing.T) {

	fmt.Println("===========================")
	fmt.Println("TestGetApiKey")
	fmt.Println("===========================")
	expectedApkiKeyValue := "12345"

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("ApiKey %s", expectedApkiKeyValue))

	apiKey, getApiKeyErr := GetAPIKey(header)

	if getApiKeyErr != nil {
		t.Errorf("unable to get API key %s\n", getApiKeyErr)
	}

	if apiKey == expectedApkiKeyValue {
		fmt.Println("PASSED!")
	} else {
		t.Errorf("expected: %s\n Recieved: %s\n", expectedApkiKeyValue, apiKey)
	}

}