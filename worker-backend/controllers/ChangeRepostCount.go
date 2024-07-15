package controllers

import (
	"context"

	"github.com/Alfazal007/worker/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) ChangeRepostCount(data Incoming) bool {
	tweetId, err := uuid.Parse(data.TweetId)
	if err != nil {
		return false
	}
	count, err := apiCfg.DB.CountRepostActual(context.Background(), tweetId)
	if data.Numb == "increment" {
		if count == 1 {
			// add new repost
			_, err := apiCfg.DB.AddNewRepost(context.Background(), database.AddNewRepostParams{
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
			_, err := apiCfg.DB.UpdateRepost(context.Background(), database.UpdateRepostParams{
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
			_, err := apiCfg.DB.DeleteRepost(context.Background(), tweetId)
			if err != nil {
				return false
			}
			return true
		}
		var newCount int32 = int32(count)
		_, err := apiCfg.DB.UpdateRepost(context.Background(), database.UpdateRepostParams{
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
