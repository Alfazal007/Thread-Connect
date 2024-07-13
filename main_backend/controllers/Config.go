package controllers

import (
	"thread-connect/internal/database"

	"github.com/cloudinary/cloudinary-go"
	"github.com/redis/go-redis/v9"
)

type ApiCfg struct {
	DB  *database.Queries
	Cld *cloudinary.Cloudinary
	Rdb *redis.Client
}
