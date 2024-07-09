package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func (apiCfg *ApiCfg) cloudinarUploader(r *http.Request, fileName string) (string, string, error) {
	// Upload a file
	uploadResult, err := apiCfg.Cld.Upload.Upload(r.Context(),
		fmt.Sprintf("uploads/%v", fileName), uploader.UploadParams{})
	if err != nil {
		fmt.Println("cloudinary error", err)
		return "", "", errors.New("error uploading to cloudinary")
	}
	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

func (apiCfg *ApiCfg) cloudinaryDeleter(r *http.Request, publicId string) error {
	_, err := apiCfg.Cld.Upload.Destroy(r.Context(), uploader.DestroyParams{PublicID: publicId})
	if err != nil {
		return err
	}
	return nil
}
