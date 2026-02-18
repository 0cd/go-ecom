package users

import (
	"log"
	"net/http"
	"strconv"

	repo "github.com/0cd/go-ecom/internal/adapters/sqlc"
	"github.com/0cd/go-ecom/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUserParams CreateUserParams
	if err := json.Read(r, &newUserParams); err != nil {
		log.Printf("Failed to parse user request: %v", err)
		http.Error(w, "invalid user request body", http.StatusBadRequest)
		return
	}

	newUser, err := h.service.CreateUser(r.Context(), newUserParams)
	if err != nil {
		log.Printf("Failed to create user in service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, newUser)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Invalid user id %s: %v", idParam, err)
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(r.Context(), id)
	if err != nil {
		log.Printf("Failed to delete user in service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusNoContent, "user deleted successfully")
}

func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		http.Error(w, "failed to retrieve users", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, users)
}

func (h *handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok || userID == 0 {
		log.Printf("Failed to extract userID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.FindUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to find user: %v", err)
		http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, user)
}

func (h *handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("query")
	if searchQuery == "" {
		log.Printf("Search query is empty")
		http.Error(w, "empty search query", http.StatusBadRequest)
		return
	}

	users, err := h.service.SearchUsers(r.Context(), pgtype.Text{
		Valid:  true,
		String: searchQuery,
	})
	if err != nil {
		log.Printf("Failed to find users (%s): %v", searchQuery, err)
		http.Error(w, "failed to search users", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, users)
}

func (h *handler) FindUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Invalid user id: %s: %v", idParam, err)
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.service.FindUserByID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to find user %d: %v", id, err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	json.Write(w, http.StatusOK, user)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Invalid user id: %s: %v", idParam, err)
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	updates := repo.UpdateUserParams{
		ID: id,
	}
	if err := json.Read(r, &updates); err != nil {
		log.Printf("Failed to parse user update request: %v", err)
		http.Error(w, "invalid user update request body", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), updates)
	if err != nil {
		log.Printf("Failed to update user in service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, updatedUser)
}

func (h *handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok || userID == 0 {
		log.Printf("Failed to extract userID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var updatePasswordParams UpdateUserPasswordParams
	if err := json.Read(r, &updatePasswordParams); err != nil {
		log.Printf("Failed to parse user password update request: %v", err)
		http.Error(w, "invalid user password update request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateUserPassword(r.Context(), userID, updatePasswordParams.OldPassword, updatePasswordParams.NewPassword)
	if err != nil {
		log.Printf("Failed to update user password in service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, "password updated successfully")
}

func (h *handler) UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok || userID == 0 {
		log.Printf("Failed to extract userID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var newEmail string
	if err := json.Read(r, &newEmail); err != nil {
		log.Printf("Failed to parse user email update request: %v", err)
		http.Error(w, "invalid user email update request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateUserEmail(r.Context(), userID, newEmail)
	if err != nil {
		log.Printf("Failed to update user password in service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, "email updated successfully")
}

func (h *handler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok || userID == 0 {
		log.Printf("Failed to extract userID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: generate an email verification token at user registration and include it in the link
	// TODO: verify it against db value here
	// TODO: implement the feature lol
	// just a static string for now

	verificationToken := r.URL.Query().Get("token")
	if verificationToken == "" {
		log.Printf("Verification token is empty")
		http.Error(w, "empty verification token", http.StatusBadRequest)
		return
	}

	err := h.service.VerifyUser(r.Context(), userID, verificationToken)
	if err != nil {
		log.Printf("Failed to verify user in service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, "user verified successfully")
}
