package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"gitlab.com/EysteinnSig/stackmap-api/pkg/users"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func validateRegisterRequest(req RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	// Add additional validation if needed (e.g., regex for email or password strength)
	return nil
}

func CreateRegisterUserHandler(userStore users.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*if r.Method != http.MethodPost {
			slog.Debug("invalid request payload", "error", err)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}*/

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Debug("invalid request payload", "error", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate the request
		if err := validateRegisterRequest(req); err != nil {
			slog.Debug("unable to validate user register request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if username or email already exists
		/*if _, err := userStore.GetUserByUsername(req.Username); err == nil {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}*/
		if _, err := userStore.GetUserByEmail(req.Email); err == nil {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}

		// Hash the password
		hashedPassword, err := HashPassword(req.Password)
		if err != nil {
			slog.Debug("unable to hash password", "error", err)
			http.Error(w, "Password not valid", http.StatusConflict)

		}

		// Create the user
		user := &users.User{
			Email:          req.Email,
			HashedPassword: hashedPassword,
		}
		if err := userStore.CreateUser(user); err != nil {
			slog.Error("error creating user", "error", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"User registered successfully"}`))
	}
}
