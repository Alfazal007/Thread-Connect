package routes

import (
	"net/http"
	"thread-connect/controllers"

	"github.com/go-chi/chi/v5"
)

func TweetRouter(apiCfg *controllers.ApiCfg) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create-tweet-media", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.CreateNewTweetWithMedia)).ServeHTTP)
	return r
}
