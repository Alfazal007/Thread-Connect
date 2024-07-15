package controllers

import (
	"context"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) GetFollowers(data Incoming) ([]string, string) {
	userId, err := uuid.Parse(data.Numb)
	if err != nil {
		var data []string
		return data, ""
	}
	username, err := apiCfg.DB.GetUsername(context.Background(), userId)
	if err != nil {
		var data []string
		return data, ""
	}

	followers, err := apiCfg.DB.GetFollowers(context.Background(), userId)
	if err != nil {
		var data []string
		return data, ""
	}
	return followers, username
}
