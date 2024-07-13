package controllers

import (
	"encoding/json"
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) RemoveRepost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue getting the user")
		return
	}
	type parameters struct {
		RepostId string `json:"repost_id"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request body")
		return
	}
	repostId, err := uuid.Parse(params.RepostId)

	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request body")
		return
	}
	repost, err := apiCfg.DB.GetRepostById(r.Context(), repostId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Repost not found")
		return
	}
	if repost.UserID != user.ID {
		helpers.RespondWithError(w, 400, "Delete your own reposts")
		return
	}
	_, err = apiCfg.DB.DeleteRepost(r.Context(), repostId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Error deleting repost")
		return
	}
	// TODO:: NEED TO ADD A PUB SUB TO DECREASE THE COUNT'
	type Message struct {
		Data string `json:"data"`
	}
	helpers.RespondWithJson(w, 201, Message{Data: "Removed the repost"})
}
