package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"gitlab.com/EysteinnSig/stackmap-api/pkg/contextdata"
	jwtutils "gitlab.com/EysteinnSig/stackmap-api/pkg/jwt"
)

const (
	claim_username_key = "userID"
)

// Middleware to validate JWT and extract userID
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("JWTMiddleware called")
		// Extract the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, `{"error": "Missing or invalid Authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
		fmt.Println("Token string: >" + tokenString + "<")
		/*_, err := jwtutils.ValidateJWT(tokenString)
		if err != nil {
			slog.Debug("Validation failed")
		} else {
			slog.Debug("Validation succeededd")

		}*/
		/*token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			secret, err := getJWTSecret()
			if err != nil {
				return "", err
			}
			return jwtSecret, nil
		})*/
		claims, err := jwtutils.ParseAccessToken(tokenString)

		//		token, err := jwtutils.ValidateJWT(tokenString)
		if err != nil { //} || !token.Valid {
			slog.Debug("Error parsing token", "error", err)
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}
		/*
			// Extract userID from token claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				fmt.Println("THIS IS SOMETHING:", claims[claim_username_key], fmt.Sprintf("%T", claims[claim_username_key]))
				userID, ok := claims[claim_username_key].(float64)
				if !ok {
					slog.Error("cannot get claims", "error", err)
					http.Error(w, `{"error": "Invalid token payload"}`, http.StatusUnauthorized)
					return
				}

				// Add userID to request context
				/*
					ctx := r.Context()
					ctx = context.WithValue(ctx, "userID", userID)*/
		/*data := getContextData(r.Context())
		if data == nil {
			data = &ContextData{}
			ctx := withContextData(r.Context(), data)
			r = r.WithContext(ctx)
		}*/
		slog.Debug("Succesfully extracted token and claims", "userID", claims.UserID)
		data, err := contextdata.GetContextData(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data.UserID = claims.UserID
		next.ServeHTTP(w, r)
		return

		//http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
	})
}
