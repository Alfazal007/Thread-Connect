package controllers

import (
	"github.com/Alfazal007/worker/internal/database"
	"github.com/Alfazal007/worker/send_mail"

	"github.com/redis/go-redis/v9"
)

type ApiCfg struct {
	DB         *database.Queries
	Rdb        *redis.Client
	GrpcClient send_mail.SendMailServiceClient
}
type Incoming struct {
	Type    string `json:"type"`
	TweetId string `json:"tweetId"`
	Numb    string `json:"numb"`
}
