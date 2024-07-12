package controllers

import (
	"encoding/json"
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) CreateRepost(w http.ResponseWriter, r *http.Request) {
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
	tweet, err := apiCfg.DB.FindATweet(r.Context(), tweetId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Tweet not found")
		return
	}
	if tweet.UserID == user.ID {
		helpers.RespondWithError(w, 400, "You cannot repost your own tweets")
		return
	}
	_, err = apiCfg.DB.GetRepost(r.Context(), database.GetRepostParams{
		UserID:  user.ID,
		TweetID: tweetId,
	})
	if err == nil {
		helpers.RespondWithError(w, 400, "Already created a repost")
		return
	}
	repost, err := apiCfg.DB.CreateNewRepost(r.Context(), database.CreateNewRepostParams{
		ID:      uuid.New(),
		TweetID: tweetId,
		UserID:  user.ID,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Error reposting")
		return
	}
	// TODO:: NEED TO ADD A PUB SUB TO INCREASE THE COUNT'
	helpers.RespondWithJson(w, 201, helpers.CustomRepostConvertor(repost))
}
