package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) RegisterUser(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	type parameter struct {
		Username string `json:"username" validate:"required,min=6,max=20"`
		Password string `json:"password" validate:"required,min=6"`
		Email    string `json:"email" validate:"required,email"`
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
			if err.Field() == "Email" {
				errMsg += "Invalid email  "
			} else if err.Field() == "Username" {
				errMsg += "Username should be greater than or equal to 6 characters and less than or equal to 15  "
			} else if err.Field() == "Password" {
				errMsg += "Password should be greater than or equal to 6 characters and less than or equal to 15  "
			}
		}
		helpers.RespondWithError(w, 400, strings.Trim(errMsg, " "))
		return
	}
	count, err := apiCfg.DB.GetUniqueUser(r.Context(), database.GetUniqueUserParams{
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue talking to the database")
		return
	}
	if count > 0 {
		helpers.RespondWithError(w, 400, "Database has same username or email")
		return
	}
	hashedPassword, err := helpers.HashPassword(params.Password)
	if err != nil {
		helpers.RespondWithError(w, 400, "Error hashing the password")
		return
	}
	// add to the database
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:       uuid.New(),
		Username: params.Username,
		Email:    params.Email,
		Password: hashedPassword,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Error talking to the database")
		return
	}
	helpers.RespondWithJson(w, 201, helpers.CustomUserReturner(user))
}
