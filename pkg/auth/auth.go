package auth

import (
	"net/http"
)

// Middleware to authorize access to a project
func authorizeProjectAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract project_id from the URL (assuming it's a query parameter for simplicity)
		projectID := r.URL.Query().Get("project_id")
		if projectID == "" {
			http.Error(w, `{"error": "Project ID is required"}`, http.StatusBadRequest)
			return
		}

		// Extract userID from the request context
		userID := r.Context().Value("userID")
		if userID == nil {
			http.Error(w, `{"error": "User ID is missing"}`, http.StatusUnauthorized)
			return
		}

		// Example authorization logic
		if !userHasAccessToProject(userID.(string), projectID) {
			http.Error(w, `{"error": "You do not have access to this project"}`, http.StatusForbidden)
			return
		}

		// If authorized, continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// Simulated function to check user access to a project
func userHasAccessToProject(userID string, projectID string) bool {
	// Replace this with real authorization logic, e.g., querying a database
	return projectID == "12345" && userID == "user1"
}

/*
// Generate JWT for demonstration purposes
func generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
*/
