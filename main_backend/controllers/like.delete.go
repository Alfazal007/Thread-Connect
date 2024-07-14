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

func (apiCfg *ApiCfg) RemoveLike(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue getting the user")
		return
	}
	type parameters struct {
		TweetId string `json:"tweet_id"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request body")
		return
	}
	tweetId, err := uuid.Parse(params.TweetId)

	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request body")
		return
	}
	_, err = apiCfg.DB.GetLike(r.Context(), database.GetLikeParams{
		UserID:  user.ID,
		TweetID: tweetId,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Like not found")
		return
	}
	_, err = apiCfg.DB.DeleteLike(r.Context(), database.DeleteLikeParams{
		UserID:  user.ID,
		TweetID: tweetId,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Error disliking")
		return
	}
	// TODO:: NEED TO ADD A PUB SUB TO DECREASE THE COUNT'
	type Like struct {
		Type    string `json:"type"`
		TweetId string `json:"tweetId"`
		Numb    string `json:"numb"`
	}
	// here type has repost type and tweetId is self explanatory and Numb shows increment or decrement
	likeType := Like{
		Type:    "like",
		TweetId: tweetId.String(),
		Numb:    "decrement",
	}
	likeTypeStr := fmt.Sprintf(`{"type":"%s","tweetId":"%s","numb":"%s"}`,
		likeType.Type, likeType.TweetId, likeType.Numb)

	err = apiCfg.Rdb.LPush(context.Background(), "worker", likeTypeStr).Err()
	if err != nil {
		println("Error adding to the redis queue")
	}
	type Message struct {
		Liked bool `json:"like"`
	}
	helpers.RespondWithJson(w, 201, Message{Liked: false})
}
