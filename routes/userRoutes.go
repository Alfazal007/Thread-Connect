package routes

import (
	"net/http"
	"thread-connect/controllers"

	"github.com/go-chi/chi/v5"
)

func UserRouter(apiCfg *controllers.ApiCfg) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", apiCfg.RegisterUser)
	r.Post("/login", apiCfg.LoginUser)
	r.Get("/current-user", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.CurrentUser)).ServeHTTP)
	r.Post("/refresh-token", controllers.VerifyRefreshJWT(apiCfg, http.HandlerFunc(apiCfg.RefreshTokens)).ServeHTTP)
	r.Put("/update-password", controllers.VerifyJWT(apiCfg, http.HandlerFunc(apiCfg.UpdatePassword)).ServeHTTP)
	return r
}
