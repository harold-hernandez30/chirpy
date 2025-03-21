package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	currentTime := time.Now().UTC()
	numericDate := jwt.NewNumericDate(currentTime)
	
	expireDate := jwt.NewNumericDate(numericDate.Add(expiresIn * time.Second))
	claim := &jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: numericDate,
		ExpiresAt: expireDate,
		Subject: userID.String(),
	}	

	claimedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, signingErr := claimedToken.SignedString([]byte(tokenSecret))

	if signingErr != nil {
		return "", signingErr
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claim := &jwt.RegisteredClaims{
		Issuer: "chirpy",
	}	


	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithLeeway(5*time.Second))
	
	if err != nil {
		return uuid.NullUUID{ Valid: false}.UUID, err
	}

	if subj, subjErr := token.Claims.GetSubject(); subjErr == nil {
		if parsedUUID, parsedUUIDErr := uuid.Parse(subj); parsedUUIDErr == nil {
			return parsedUUID, parsedUUIDErr
		}
	}

	return uuid.NullUUID{ Valid: false}.UUID, err
}