package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"thread-connect/helpers"
	"thread-connect/internal/database"
	"time"

	"github.com/go-playground/validator/v10"
)

func (apiCfg *ApiCfg) LoginUser(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	type parameter struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}
	var params parameter
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	err = validate.Struct(params)
	if err != nil {
		var errMsg string
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Username" {
				errMsg += "Username should be greater than or equal to 6 characters and less than or equal to 15  "
			} else if err.Field() == "Password" {
				errMsg += "Password should be greater than or equal to 6 characters and less than or equal to 15  "
			}
		}
		helpers.RespondWithError(w, 400, strings.Trim(errMsg, " "))
		return
	}
	user, err := apiCfg.DB.GetUserByName(r.Context(), params.Username)
	if err != nil {
		helpers.RespondWithError(w, 400, "User not found")
		return
	}
	isPasswordValid := helpers.CompareHash(params.Password, user.Password)
	if !isPasswordValid {
		helpers.RespondWithError(w, 400, "Incorrect password")
		return
	}
	accessToken, err := GenerateJWT(user)
	if err != nil {
		helpers.RespondWithError(w, 400, "Error generating access token")
		return
	}
	refreshToken, err := GenerateRefreshToken(user)
	if err != nil {
		helpers.RespondWithError(w, 400, "Error generating refresh token")
		return
	}
	user, err = apiCfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		RefreshToken: sql.NullString{String: refreshToken, Valid: true},
		ID:           user.ID,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Error updating refresh token in the database")
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
