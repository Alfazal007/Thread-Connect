package routes

import (
	"thread-connect/controllers"

	"github.com/go-chi/chi/v5"
)

func UserRouter(apiCfg *controllers.ApiCfg) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", apiCfg.RegisterUser)
	return r
}
