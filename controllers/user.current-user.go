package controllers

import (
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"
)

func (apiCfg *ApiCfg) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	customUser := helpers.CustomUserReturner(user)
	helpers.RespondWithJson(w, 200, customUser)
}
