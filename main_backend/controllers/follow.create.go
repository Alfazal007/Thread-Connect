package controllers

import (
	"encoding/json"
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) FollowUser(w http.ResponseWriter, r *http.Request) {
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
		helpers.RespondWithError(w, 400, "Incorrect request body")
		return
	}
	followingId, err := uuid.Parse(params.Following)
	if err != nil {
		helpers.RespondWithError(w, 400, "Incorrect request body")
		return
	}
	if followingId == user.ID {
		helpers.RespondWithError(w, 400, "Cannot follow yourself")
		return
	}
	// check if already a relationship is present
	_, err = apiCfg.DB.AlreadyFollowing(r.Context(), database.AlreadyFollowingParams{
		Follower:  user.ID,
		Following: followingId,
	})
	if err == nil {
		helpers.RespondWithError(w, 400, "Already following")
		return
	}
	// add the following
	_, err = apiCfg.DB.CreateNewFollow(r.Context(), database.CreateNewFollowParams{
		ID:        uuid.New(),
		Follower:  user.ID,
		Following: followingId,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Failed to follow")
		return
	}
	type Data struct {
		Message string `json:"message"`
	}
	helpers.RespondWithJson(w, 200, Data{Message: "Following the user"})
}
