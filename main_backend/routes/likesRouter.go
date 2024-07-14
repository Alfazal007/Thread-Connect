package routes

import (
	"net/http"
	"thread-connect/controllers"

	"github.com/go-chi/chi/v5"
)

func LikeRouter(apiCfg *controllers.ApiCfg) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/create", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.CreateLike)).ServeHTTP)
	r.Delete("/delete", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.RemoveLike)).ServeHTTP)
	return r
}
