package controllers

import (
	"thread-connect/internal/database"

	"github.com/cloudinary/cloudinary-go"
)

type ApiCfg struct {
	DB  *database.Queries
	Cld *cloudinary.Cloudinary
}
