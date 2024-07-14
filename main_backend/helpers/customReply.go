package helpers

import (
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

type CustomReply struct {
	ID           uuid.UUID `json:"tweet_id"`
	Content      string    `json:"content"`
	Media        string    `json:"media"`
	UserID       uuid.UUID `json:"user_id"`
	CreatedAt    string    `json:"created_at"`
	ReplyTweetID uuid.UUID `json:"reply_tweet_id"`
}

func CustomReplyConvertor(tweet database.Tweet) CustomReply {
	return CustomReply{
		ID:           tweet.ID,
		Content:      tweet.Content.String,
		Media:        tweet.Media.String,
		UserID:       tweet.UserID,
		CreatedAt:    tweet.CreatedAt.String(),
		ReplyTweetID: tweet.ReplyTweetID.UUID,
	}
}
