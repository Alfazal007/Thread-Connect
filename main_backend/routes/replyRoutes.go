package routes

import (
	"net/http"
	"thread-connect/controllers"

	"github.com/go-chi/chi/v5"
)

func ReplyRouter(apiCfg *controllers.ApiCfg) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create-reply", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.ReplyTweet)).ServeHTTP)
	r.Delete("/delete-reply", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.ReplyDelete)).ServeHTTP)
	return r
}
