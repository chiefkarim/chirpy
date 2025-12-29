package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/chiefkarim/chirpy/internal/auth"
	"github.com/chiefkarim/chirpy/internal/database"
)

func (config *apiConfig) RevokeRefreshToken(response http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithJSON(response, http.StatusUnauthorized, map[string]string{"error": "Please provide valid access token in the header."})
		return
	}
	_, err = config.dbQueries.ExpireToken(request.Context(), database.ExpireTokenParams{
		Token: token,
		RevokedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})
	respondWithJSON(response, http.StatusNoContent, nil)
}
