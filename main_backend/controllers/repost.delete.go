package controllers

import (
	"context"
	"encoding/json"
	"fmt"
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
	type Repost struct {
		Type    string `json:"type"`
		TweetId string `json:"tweetId"`
		Numb    string `json:"numb"`
	}
	// here type has repost type and tweetId is self explanatory and Numb shows increment or decrement
	repostType := Repost{
		Type:    "repost",
		TweetId: repost.TweetID.String(),
		Numb:    "decrement",
	}
	repostTypeStr := fmt.Sprintf(`{"type":"%s","tweetId":"%s","numb":"%s"}`,
		repostType.Type, repostType.TweetId, repostType.Numb)

	err = apiCfg.Rdb.LPush(context.Background(), "worker", repostTypeStr).Err()
	if err != nil {
		println("Error adding to the redis queue")
	}
	type Message struct {
		Data string `json:"data"`
	}
	helpers.RespondWithJson(w, 201, Message{Data: "Removed the repost"})
}
