package middleware

import (
	"net/http"

	"github.com/0cd/go-ecom/internal/users"
)

func AdminMiddleware(userService users.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("userID").(int64)
			if !ok || userID == 0 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userService.FindUserByID(r.Context(), userID)
			if err != nil {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			if !user.IsAdmin {
				http.Error(w, "forbidden: admin only", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
