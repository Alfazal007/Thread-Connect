package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"thread-connect/helpers"
	"thread-connect/internal/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(user database.User) (string, error) {
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("DB url is not found in env variables")
	}
	secretKey := []byte(jwtSecret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["username"] = user.Username

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(user database.User) (string, error) {
	jwtSecret := os.Getenv("SECRET_KEY_REFRESH")
	if jwtSecret == "" {
		log.Fatal("DB url is not found in env variables")
	}
	secretKey := []byte(jwtSecret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(240 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(apiCfg *ApiCfg, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access-token")
		var jwtToken string
		if err != nil {
			if err == http.ErrNoCookie {
				authorization := r.Header.Get("Authorization")
				if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
					http.Error(w, "Authorization header missing or improperly formatted", http.StatusUnauthorized)
					helpers.RespondWithError(w, 400, "No headers provided")
					return
				}

				jwtToken = strings.TrimPrefix(authorization, "Bearer ")
			} else {
				helpers.RespondWithError(w, 400, "Error reading cookie, try logging in again")
				return
			}
		} else {
			jwtToken = cookie.Value
		}
		// Verify the JWT token
		jwtSecret := os.Getenv("SECRET_KEY")
		if jwtSecret == "" {
			helpers.RespondWithError(w, 400, "Server error")
			return
		}
		if jwtToken == "" {
			helpers.RespondWithError(w, 400, "Provide cookie")
			return
		}
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			helpers.RespondWithError(w, 401, fmt.Sprintf("Invalid token here %v", err))
			return
		}
		if !token.Valid {
			helpers.RespondWithError(w, 401, fmt.Sprintf("Invalid token %v", err))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.RespondWithError(w, 400, "Invalid claims login again")
			return
		}

		username := claims["username"].(string)
		id := claims["user_id"].(string)

		user, err := apiCfg.DB.GetUserByName(r.Context(), username)
		if err != nil {
			helpers.RespondWithError(w, 400, "Some manpulation done with the token")
			return
		}
		idUUID, err := uuid.Parse(id)
		if err != nil {
			helpers.RespondWithError(w, 400, "Some manpulation done with the token")
			return
		}
		if idUUID != user.ID {
			helpers.RespondWithError(w, 400, "Some manipulations done with the token try again")
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
