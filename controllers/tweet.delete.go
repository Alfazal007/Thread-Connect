package controllers

import (
	"encoding/json"
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) DeleteTweet(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	type parameter struct {
		TweetId string `json:"tweetId"`
	}
	var params parameter
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Error getting response body")
		return
	}
	tweetId, err := uuid.Parse(params.TweetId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid tweet id")
		return
	}
	tweet, err := apiCfg.DB.FindATweet(r.Context(), tweetId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Not found this tweet in the database")
		return
	}
	if tweet.UserID != user.ID {
		helpers.RespondWithError(w, 400, "You are not this tweet author")
		return
	}

	if tweet.Media.String != "" {
		err = apiCfg.cloudinaryDeleter(r, tweet.PublicID.String)
		if err != nil {
			helpers.RespondWithError(w, 400, "Error deleting the media")
			return
		}
	}
	_, err = apiCfg.DB.DeleteATweet(r.Context(), database.DeleteATweetParams{
		ID:     tweetId,
		UserID: user.ID,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue deleting the tweet")
		return
	}
	type Data struct {
		Message string `json:"message"`
	}
	helpers.RespondWithJson(w, 200, Data{Message: "Deleted tweet successfully"})
}
