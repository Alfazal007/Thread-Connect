package routes

import (
	"net/http"
	"thread-connect/controllers"

	"github.com/go-chi/chi/v5"
)

func FollowRouter(apiCfg *controllers.ApiCfg) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/follow", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.FollowUser)).ServeHTTP)
	r.Post("/unfollow", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.UnfollowUser)).ServeHTTP)
	return r
}
