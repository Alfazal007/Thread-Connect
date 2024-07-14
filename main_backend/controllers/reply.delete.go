package controllers

import (
	"encoding/json"
	"net/http"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) ReplyDelete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "User not found")
		return
	}
	type parameters struct {
		ReplyId string `json:"replyId"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Error decoding the data")
		return
	}

	replyId, err := uuid.Parse(params.ReplyId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid reply id")
		return
	}

	_, err = apiCfg.DB.DeleteReply(r.Context(), database.DeleteReplyParams{
		ID:     replyId,
		UserID: user.ID,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, "Error deleting the reply")
		return
	}
	type Data struct {
		Message string `json:"message"`
	}
	helpers.RespondWithJson(w, 200, Data{Message: "Deleted the reply"})
}
