package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTimeDiff(t *testing.T) {
	

	fmt.Println("===========================")
	fmt.Println("TestTimeDiff")
	fmt.Println("===========================")
	location, err := time.LoadLocation("UTC")

	if err != nil {
		t.Errorf("could not load location: %s\n", err)
	}
	
	refDate := time.Date(2000, time.August, 9, 10, 0, 0, 0, location)

	threeSecondsLater := refDate.Add(3 * time.Second)

	compareTime(refDate, threeSecondsLater)
	if threeSecondsLater.Second() == 3 {
		fmt.Println("3 seconds added")
	} else {
		t.Errorf("failed adding 3 seconds")
	}

	hoursLater := refDate.Add(5 * time.Hour)


	compareTime(refDate, hoursLater)
	if hoursLater.Hour() == 15 {
		fmt.Println("5 hours added")
	} else {
		t.Errorf("failed adding 5 hours")
	}
}

func TestJWTExpiry(t *testing.T) {
	baseUUID, err := uuid.NewRandom()
	testSecret := "someSecret"

	if err != nil {
		t.Errorf("unable to create UUID")
	}
	// 1hr after - should be valid
	oneHourToken, oneHourTokenErr := MakeJWT(baseUUID, testSecret, 1 * time.Hour)
	if err := makeJWTTestHelper(oneHourTokenErr, t); err != nil {
		return
	}

	if err := validateJWTTestHelper(oneHourToken, testSecret, true, t); err != nil {
		return
	}


	// 1hr before - should be invalid
	oneHourBeforeToken, oneHourBeforeTokenErr := MakeJWT(baseUUID, testSecret, -1 * time.Hour)
	if err := makeJWTTestHelper(oneHourBeforeTokenErr, t); err != nil {
		return
	}

	if err := validateJWTTestHelper(oneHourBeforeToken, testSecret, false, t); err != nil {
		return
	}

	// 3 secs after - should be valid
	within5SecLeeway, within5SecLeewayErr := MakeJWT(baseUUID, testSecret, 3 * time.Second)
	if err := makeJWTTestHelper(within5SecLeewayErr, t); err != nil {
		return
	}

	if err := validateJWTTestHelper(within5SecLeeway, testSecret, true, t); err != nil {
		return
	}

	// 10 secs after - should be invalid
	outside5SecLeeway, outside5SecLeewayErr := MakeJWT(baseUUID, testSecret, 10 * time.Second)
	if err := makeJWTTestHelper(outside5SecLeewayErr, t); err != nil {
		return
	}

	if err := validateJWTTestHelper(outside5SecLeeway, testSecret, false, t); err != nil {
		return
	}
}

func makeJWTTestHelper(makeJWTErr error, t *testing.T) error{
	if makeJWTErr != nil {
		t.Errorf("Unable to create onehourToken: %s", makeJWTErr)
	}
	return makeJWTErr
}

func validateJWTTestHelper(tokenString, testSecret string, shouldBeValid bool, t *testing.T) error{
		_, tokenErr := ValidateJWT(tokenString, testSecret)
	
	if shouldBeValid && tokenErr != nil {
		t.Errorf("Failed: %s\n", tokenErr)
		return fmt.Errorf("expecting valid token but is invalid")
	}

	if !shouldBeValid && tokenErr == nil {
		t.Errorf("Failed: %s\n", tokenErr)
		return fmt.Errorf("expecting invalid token but is valid")
	}

	fmt.Println("PASS!")
	return nil
}

func compareTime(refTime time.Time, other time.Time) {
	fmt.Println("==================================")
	fmt.Printf("\tRefTime\t\tOtherTime\n")
	fmt.Printf("Year:\t%d\t\t%d\n", refTime.Year(), other.Year())
	fmt.Printf("Month:\t%d\t\t%d\n", refTime.Month(), other.Month())
	fmt.Printf("Day:\t%d\t\t%d\n", refTime.Day(), other.Day())
	fmt.Printf("Hour:\t%d\t\t%d\n", refTime.Hour(), other.Hour())
	fmt.Printf("Sec:\t%d\t\t%d\n", refTime.Second(), other.Second())
}

func TestMakeJWT(t *testing.T) {
	fmt.Println("")
	fmt.Println("===========================")
	fmt.Println("TestMakeJWT")
	fmt.Println("===========================")
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