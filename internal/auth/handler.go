package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/0cd/go-ecom/internal/env"
	"github.com/0cd/go-ecom/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// TODO: remove password_hash from the return
func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var registerParams userLoginAndRegisterParams
	if err := json.Read(r, &registerParams); err != nil {
		log.Printf("Failed to parse user request: %v", err)
		http.Error(w, "invalid user request body", http.StatusBadRequest)
		return
	}

	newUser, err := h.service.Register(r.Context(), registerParams.Email, registerParams.Password)
	if err != nil {
		log.Printf("Failed to register user in service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, newUser)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginParams userLoginAndRegisterParams
	if err := json.Read(r, &loginParams); err != nil {
		log.Printf("Failed to parse user request: %v", err)
		http.Error(w, "invalid user request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), loginParams.Email, loginParams.Password)
	if err != nil {
		log.Printf("Failed to login in service: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := h.service.GenerateTokens(user.ID)
	if err != nil {
		log.Printf("Failed to generate tokens: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	setAuthCookies(w, accessToken, refreshToken)
	json.Write(w, http.StatusOK, "logged in successfully")
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	clearAuthCookies(w)
	json.Write(w, http.StatusOK, "logged out successfully")
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Printf("Failed to extract cookie: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := ParseRefreshToken(cookie.Value)
	if err != nil {
		log.Printf("Failed to parse refresh token: %v", err)
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, _, err := h.service.GenerateTokens(claims.UserID)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		http.Error(w, "something went wrong", http.StatusUnauthorized)
		return
	}

	setAuthCookies(w, accessToken, cookie.Value)
	json.Write(w, http.StatusOK, "access token refreshed")
}

// TODO: find a better solution for secure bool
func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   !(env.GetString("ENVIRONMENT", "prod") == "dev"),
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(10 * time.Minute),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   !(env.GetString("ENVIRONMENT", "prod") == "dev"),
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}

func clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   !(env.GetString("ENVIRONMENT", "prod") == "dev"),
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   !(env.GetString("ENVIRONMENT", "prod") == "dev"),
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})
}
