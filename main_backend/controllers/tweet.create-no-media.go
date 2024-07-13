package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) TweetNoMedia(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "User not found")
		return
	}
	type parameters struct {
		Content string `json:"content"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Error decoding the data")
		return
	}
	if len(params.Content) < 1 {
		helpers.RespondWithError(w, 400, "No content given")
		return
	}
	tweet, err := apiCfg.DB.CreateNewTweetNoMedia(r.Context(), database.CreateNewTweetNoMediaParams{
		ID:        uuid.New(),
		Content:   sql.NullString{String: params.Content, Valid: true},
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Error creating the tweet")
		return
	}
	helpers.RespondWithJson(w, 201, helpers.CustomTweetConvertor(tweet))
}
