package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	jwtutils "gitlab.com/EysteinnSig/stackmap-api/pkg/jwt"
)

func CreateRefreshTokenHandler(service RefreshTokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
		/*if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}*/

		type RefreshRequest struct {
			RefreshToken string `json:"refresh_token"`
		}

		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Debug("Unable to decode body", "error", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		slog.Debug("Refreshing token", "old token", req.RefreshToken)

		//jwtutils.ParseRefreshToken(req.RefreshToken)
		userID, err := service.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			slog.Debug("Unable to validate refresh token", "error", err)
			http.Error(w, `{"error": "Invalid or expired refresh token"}`, http.StatusUnauthorized)
			return
		}
		/*

			token, err := jwtutils.ValidateJWT(req.RefreshToken) // validateRefreshToken(req.RefreshToken)
			if err != nil || !token.Valid {
				http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}*/

		/*if err != nil {
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
			return
		}*/
		if err := service.InvalidateRefreshToken(req.RefreshToken); err != nil {
			slog.Error("Failed to invalidate refresh token", "token", req.RefreshToken, "error", err)
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		// Generate new tokens
		newRefreshToken, err := service.GenerateRefreshToken(userID)
		if err != nil {
			slog.Error("Failed to generate refresh token", "error", err)
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		newAccessToken, err := jwtutils.NewAccessToken(jwtutils.InitUserClaims(userID))
		//refreshToken, err := generateRefreshToken(userID)
		if err != nil {
			http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//json.NewEncoder(w).Encode(map[string]string{"token": token})
		json.NewEncoder(w).Encode(map[string]string{"access_token": newAccessToken, "refresh_token": newRefreshToken})
	}
}

/*
func RefreshTokenHandler(service users.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type RefreshRequest struct {
			RefreshToken string `json:"refresh_token"`
		}

		var body RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			slog.Debug("invalid request", "error", err)
			http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
			return
		}

		//token, err := jwtutils.ValidateJWT(tokenString)
		fmt.Println("TOKEN: ", body.RefreshToken)

		userClaimsParsed, err := jwtutils.ParseAccessToken(body.RefreshToken)
		if err != nil || userClaimsParsed.Valid() != nil {
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		slog.Info("successful refreshed token")
		// Respond with token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//json.NewEncoder(w).Encode(map[string]string{"token": token})
		json.NewEncoder(w).Encode(map[string]string{"access_token": signedAccessToken, "refresh_token": signedRefreshToken})
	}
}
*/
