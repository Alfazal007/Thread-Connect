package controllers

import (
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"thread-connect/helpers"
	"thread-connect/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiCfg) CreateNewTweetWithMedia(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		helpers.RespondWithError(w, 400, "Issue with finding the user from the database")
		return
	}
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		helpers.RespondWithError(w, 400, "Media sent is too big")
		return
	}
	content := r.FormValue("content")
	if content == "" {
		helpers.RespondWithError(w, 400, "Content should not be empty")
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("media")
	if err != nil {
		helpers.RespondWithError(w, 400, "Error retreiving the file")
		return
	}
	defer file.Close()
	uploadDir := "uploads"

	// Generate a random integer between 0 and 99.
	randomNumber := rand.Intn(1000)

	dstPath := filepath.Join(uploadDir, fmt.Sprintf("%v%v", randomNumber, handler.Filename))
	dst, err := os.Create(dstPath)
	if err != nil {
		os.Remove(dstPath)
		helpers.RespondWithError(w, 400, "Error saving the file")
		return
	}
	defer dst.Close()

	// Copy the uploaded file's content to the new file
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(dstPath)
		helpers.RespondWithError(w, 400, "Error saving the file")
		return
	}
	url, publicId, err := apiCfg.cloudinarUploader(r, fmt.Sprintf("%v%v", randomNumber, handler.Filename))
	if err != nil {
		os.Remove(dstPath)
		helpers.RespondWithError(w, 400, "Error uploading to cloudinary")
		return
	}
	os.Remove(dstPath)
	uploadedTweet, err := apiCfg.DB.CreateNewTweet(r.Context(), database.CreateNewTweetParams{
		ID:       uuid.New(),
		Content:  sql.NullString{String: content, Valid: true},
		Media:    sql.NullString{String: url, Valid: true},
		PublicID: sql.NullString{String: publicId, Valid: true},
		UserID:   user.ID,
	})
	if err != nil {
		os.Remove(dstPath)
		helpers.RespondWithError(w, 400, "Error writing to the database, try again later")
		return
	}
	helpers.RespondWithJson(w, 200, helpers.CustomTweetConvertor(uploadedTweet))
}
