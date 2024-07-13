package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"thread-connect/helpers"
	"thread-connect/internal/database"
	"time"

	"github.com/go-playground/validator/v10"
)

func (apiCfg *ApiCfg) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "User not found")
		return
	}
	validate := validator.New()

	type parameter struct {
		OldPassword string `json:"oldpassword" validate:"required,min=6"`
		NewPassword string `json:"newpassword" validate:"required,min=6"`
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
			if err.Field() == "OldPassword" {
				errMsg += "Old password should be greater than or equal to 6 characters and less than or equal to 15  "
			} else if err.Field() == "NewPassword" {
				errMsg += "New Password should be greater than or equal to 6 characters and less than or equal to 15  "
			}
		}
		helpers.RespondWithError(w, 400, strings.Trim(errMsg, " "))
		return
	}
	// compare password
	isValidPassword := helpers.CompareHash(params.OldPassword, user.Password)

	if !isValidPassword {
		helpers.RespondWithError(w, 400, "Incorrect old password")
		return
	}
	newPassword, err := helpers.HashPassword(params.NewPassword)
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue hashing the password")
		return
	}
	updatedUser, err := apiCfg.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Password:  newPassword,
		ID:        user.ID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue updating the password in the database")
		return
	}
	helpers.RespondWithJson(w, 200, helpers.CustomUserReturner(updatedUser))
}
