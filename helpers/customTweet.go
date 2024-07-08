package helpers

import (
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

type CustomTweet struct {
	ID      uuid.UUID `json:"tweet_id"`
	Content string    `json:"content"`
	Media   string    `json:"media"`
	UserID  uuid.UUID `json:"user_id"`
}

func CustomTweetConvertor(tweet database.Tweet) CustomTweet {
	return CustomTweet{
		ID:      tweet.ID,
		Content: tweet.Content.String,
		Media:   tweet.Media.String,
		UserID:  tweet.UserID,
	}
}
