package main

import (
	"fmt"
	"net/http"

	"github.com/harold-hernandez30/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevokeRefreshToken(res http.ResponseWriter, req *http.Request) {

	refreshToken, getBearerTokenErr := auth.GetBearerToken(req.Header)

	if getBearerTokenErr != nil {
		fmt.Printf("could not get refresh token from heading: %s\n", getBearerTokenErr)
		res.WriteHeader(500)
		return
	}

	userAndRefreshToken, getUserFromRefreshTokenErr := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)

	if getUserFromRefreshTokenErr != nil {
		fmt.Printf("could not get user from refresh token %s\n", getUserFromRefreshTokenErr)
		res.WriteHeader(500)
		return
	}

	if revokeRefreshTokenErr := cfg.db.RevokeRefreshToken(req.Context(), userAndRefreshToken.RefreshToken.Token); revokeRefreshTokenErr == nil {
		res.WriteHeader(http.StatusNoContent)
	} else {
		
		fmt.Printf("could not revoke refresh token %s\n", getUserFromRefreshTokenErr)
		res.WriteHeader(500)
	}

	
}