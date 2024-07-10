package routes

import (
	"net/http"
	"thread-connect/controllers"

	"github.com/go-chi/chi/v5"
)

func TweetRouter(apiCfg *controllers.ApiCfg) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create-tweet-media", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.CreateNewTweetWithMedia)).ServeHTTP)
	r.Delete("/delete-tweet", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.DeleteTweet)).ServeHTTP)
	r.Post("/create-tweet-no-media", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.TweetNoMedia)).ServeHTTP)
	return r
}
