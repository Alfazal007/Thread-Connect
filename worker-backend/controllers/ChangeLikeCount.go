package controllers

import (
	"context"

	"github.com/Alfazal007/worker/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) ChangeLikeCount(data Incoming) bool {
	tweetId, err := uuid.Parse(data.TweetId)
	if err != nil {
		return false
	}
	count, err := apiCfg.DB.CountLikesActual(context.Background(), tweetId)
	if data.Numb == "increment" {
		if count == 1 {
			// add new repost
			_, err := apiCfg.DB.AddNewLike(context.Background(), database.AddNewLikeParams{
				TweetID: tweetId,
				Count:   1,
			})
			if err != nil {
				return false
			} else {
				return true
			}
		} else {
			var newCount int32 = int32(count)
			_, err := apiCfg.DB.UpdateLike(context.Background(), database.UpdateLikeParams{
				TweetID: tweetId,
				Count:   newCount,
			})
			if err != nil {
				return false
			} else {
				return true
			}
		}
	} else {
		if count == 0 {
			// delete the row from repost_count
			_, err := apiCfg.DB.DeleteLike(context.Background(), tweetId)
			if err != nil {
				return false
			}
			return true
		}
		var newCount int32 = int32(count)
		_, err := apiCfg.DB.UpdateLike(context.Background(), database.UpdateLikeParams{
			TweetID: tweetId,
			Count:   newCount,
		})
		if err != nil {
			return false
		} else {
			return true
		}
	}
}
