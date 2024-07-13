package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"
	"time"
)

func (apiCfg *ApiCfg) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	type params struct {
		RefreshToken string `json:"refresh-token"`
	}
	var parameters params
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue decoding the data")
		return
	}
	if parameters.RefreshToken != user.RefreshToken.String {
		helpers.RespondWithError(w, 400, "Incorrect refresh token")
		return
	}
	// update tokens
	accessToken, err := GenerateJWT(user)
	if err != nil {
		helpers.RespondWithError(w, 400, "Could not generate access token")
		return
	}
	refreshToken, err := GenerateRefreshToken(user)
	if err != nil {
		helpers.RespondWithError(w, 400, "Could not generate refresh token")
		return
	}
	_, err = apiCfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		RefreshToken: sql.NullString{String: refreshToken, Valid: true},
		ID:           user.ID,
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Could not update the database")
		return
	}
	cookie1 := http.Cookie{
		Name:     "access-token",
		Value:    accessToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, &cookie1)
	cookie2 := http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 240),
	}
	http.SetCookie(w, &cookie2)
	type Tokens struct {
		AccessToken  string `json:"access-token"`
		RefreshToken string `json:"refresh-token"`
	}
	helpers.RespondWithJson(w, 200, Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
