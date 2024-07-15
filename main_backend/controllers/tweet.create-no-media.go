package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
	type Notification struct {
		Type    string `json:"type"`
		TweetId string `json:"tweetId"`
		Numb    string `json:"numb"`
	}
	notificationType := Notification{
		Type:    "notification",
		TweetId: tweet.ID.String(),
		Numb:    user.ID.String(),
	}
	notificationTypeStr := fmt.Sprintf(`{"type":"%s","tweetId":"%s","numb":"%s"}`,
		notificationType.Type, notificationType.TweetId, notificationType.Numb)

	err = apiCfg.Rdb.LPush(context.Background(), "worker", notificationTypeStr).Err()
	if err != nil {
		println("Error adding to the redis queue")
	}
	helpers.RespondWithJson(w, 201, helpers.CustomTweetConvertor(tweet))
}
