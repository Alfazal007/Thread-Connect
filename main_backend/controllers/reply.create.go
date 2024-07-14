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

func (apiCfg *ApiCfg) ReplyTweet(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "User not found")
		return
	}
	type parameters struct {
		Content string `json:"content"`
		TweetId string `json:"tweetId"`
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
	tweetId, err := uuid.Parse(params.TweetId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid tweet id")
		return
	}
	tweet, err := apiCfg.DB.CreateNewReply(r.Context(), database.CreateNewReplyParams{
		ID:           uuid.New(),
		Content:      sql.NullString{String: params.Content, Valid: true},
		UserID:       user.ID,
		CreatedAt:    time.Now().UTC(),
		ReplyTweetID: uuid.NullUUID{UUID: tweetId, Valid: true},
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Error creating the reply")
		return
	}
	helpers.RespondWithJson(w, 201, helpers.CustomReplyConvertor(tweet))
}
