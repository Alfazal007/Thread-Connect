package controllers

import (
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) GetTweet(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	tweetIdStr := chi.URLParam(r, "tweetId")
	tweetId, err := uuid.Parse(tweetIdStr)
	if err != nil {
		helpers.RespondWithError(w, 400, "Incorrect tweet id")
		return
	}
	tweet, err := apiCfg.DB.FindATweet(r.Context(), tweetId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Tweet not found")
		return
	}
	helpers.RespondWithJson(w, 200, helpers.CustomTweetConvertor(tweet))
}
