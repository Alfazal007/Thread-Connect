package helpers

import (
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

type CustomReposter struct {
	ID      uuid.UUID `json:"repost_id"`
	TweetId uuid.UUID `json:"tweet_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func CustomRepostConvertor(repost database.Repost) CustomReposter {
	return CustomReposter{
		ID:      repost.ID,
		TweetId: repost.TweetID,
		UserID:  repost.UserID,
	}
}
