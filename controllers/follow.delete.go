package controllers

import (
	"encoding/json"
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "User not found")
		return
	}
	type parameters struct {
		Following string `json:"following"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request body")
		return
	}
	followingId, err := uuid.Parse(params.Following)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request body")
		return
	}
	_, err = apiCfg.DB.AlreadyFollowing(r.Context(), database.AlreadyFollowingParams{
		Follower:  user.ID,
		Following: followingId,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Not following")
		return
	}
	_, err = apiCfg.DB.UnfollowUser(r.Context(), database.UnfollowUserParams{
		Follower:  user.ID,
		Following: followingId,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Issue unfollowing the user, try again later")
		return
	}
	type Message struct {
		Data string `json:"data"`
	}
	helpers.RespondWithJson(w, 200, Message{Data: "Unfollow successful"})
}
