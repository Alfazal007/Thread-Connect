package controllers

import (
	"worker/internal/database"

	"github.com/redis/go-redis/v9"
)

type ApiCfg struct {
	DB  *database.Queries
	Rdb *redis.Client
}
type Incoming struct {
	Type    string `json:"type"`
	TweetId string `json:"tweetId"`
	Numb    string `json:"numb"`
}
