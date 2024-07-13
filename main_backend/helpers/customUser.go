package helpers

import (
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

type CustomUser struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Id           uuid.UUID `json:"id"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func CustomUserReturner(user database.User) CustomUser {
	return CustomUser{
		Username:     user.Username,
		Email:        user.Email,
		Id:           user.ID,
		RefreshToken: user.RefreshToken.String,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}
}