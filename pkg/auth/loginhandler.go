package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	jwtutils "gitlab.com/EysteinnSig/stackmap-api/pkg/jwt"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/users"
	"golang.org/x/crypto/bcrypt"
)

func CreateLoginHandler(userstore users.UserStore, tokenstore RefreshTokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var body LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			slog.Debug("invalid request", "error", err)
			http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
			return
		}

		// Fetch user from the database
		user, err := userstore.GetUserByEmail(body.Email)
		if err != nil {
			slog.Debug("unable to get user", "error", err)
			http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// Check password
		if err := CheckPassword(user.HashedPassword, body.Password); err != nil {
			slog.Debug("password check failed", "error", err)
			http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		accessToken, err := jwtutils.NewAccessToken(jwtutils.InitUserClaims(user.ID))
		if err != nil {
			slog.Debug("could not generate access token", "error", err)
			http.Error(w, `{"error": "Could not generate token"}`, http.StatusInternalServerError)
			return
		}

		refreshToken, err := tokenstore.GenerateRefreshToken(user.ID)
		if err != nil {
			slog.Debug("could not generate refresh token", "error", err)
			http.Error(w, `{"error": "Could not generate refresh token"}`, http.StatusInternalServerError)
		}

		/*userclaims := jwtutils.UserClaims{
			UserID: user.ID,
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			},
		}
		signedAccessToken, err := jwtutils.NewAccessToken(userclaims)
		if err != nil {
			slog.Debug("could not generate access token", "error", err)
			http.Error(w, `{"error": "Could not generate token"}`, http.StatusInternalServerError)
			return
		}

		// Should not contain any data
		refreshClaims := jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		}

		signedRefreshToken, err := jwtutils.NewRefreshToken(refreshClaims)
		if err != nil {
			slog.Debug("could not generate refresh token", "error", err)
			http.Error(w, `{"error": "Could not generate token"}`, http.StatusInternalServerError)
		}*/

		// Generate JWT
		/*token, err := jwtutils.GenerateJWT(user.ID)
		if err != nil {
			slog.Debug("could not generate token", "error", err)
			http.Error(w, `{"error": "Could not generate token"}`, http.StatusInternalServerError)
			return
		}*/

		slog.Info("successful login", "user", user.Email)
		// Respond with token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//json.NewEncoder(w).Encode(map[string]string{"token": token})
		//json.NewEncoder(w).Encode(map[string]string{"access_token": signedAccessToken, "refresh_token": signedRefreshToken})
		json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken, "refresh_token": refreshToken})

	}
}

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword verifies a bcrypt hashed password against a plaintext password.
func CheckPassword(hashedPassword, password string) error {
	//return nil
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

/*
import (
	"encoding/json"
	"net/http"

	"gitlab.com/EysteinnSig/stackmap-api/pkg/users"
	"golang.org/x/crypto/bcrypt"
)

func createLoginHandler(service users.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var body LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
			return
		}

		// Fetch user from the database
		user, err := service.GetUserByName(body.Username)
		if err != nil {
			http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// Check password
		if err := CheckPassword(user.HashedPassword, body.Password); err != nil {
			http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// Generate JWT
		token, err := generateJWT(user.Name)
		if err != nil {
			http.Error(w, `{"error": "Could not generate token"}`, http.StatusInternalServerError)
			return
		}

		// Respond with token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword verifies a bcrypt hashed password against a plaintext password.
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
*/
