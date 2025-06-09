package projects

import (
	"log/slog"
	"net/http"

	"gitlab.com/EysteinnSig/stackmap-api/pkg/contextdata"
)

// Middleware to authorize access to a project
func CreateAuthProjectAccessMiddleware(projectStore ProjectStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("Verifying project access")
			// Extract project_id from the URL (assuming it's a query parameter for simplicity)
			//projectID := r.URL.Query().Get("project_name")
			projectID := r.PathValue("project_name")
			if projectID == "" {
				slog.Error("Project ID is required")
				http.Error(w, `{"error": "Project ID is required"}`, http.StatusBadRequest)
				return
			}

			// Extract userID from the request context
			cdata, err := contextdata.GetContextData(r.Context())
			if err != nil {
				slog.Error("Failed to get context data", "error", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			userID := cdata.UserID

			// Example authorization logic
			/*projects, err := projectStore.GetProjectsOwnedBy([]int{userID})
			if err != nil {
				slog.Error("Failed to get projects owned by user", "error", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			for _, project := range projects {
				if project.Name == projectID {
					next.ServeHTTP(w, r)
					return
				}
			}*/
			if userHasAccessToProject(userID, projectID, projectStore) {
				slog.Debug("Authorized access", "userID", userID, "projectID", projectID)
				next.ServeHTTP(w, r)
				return
			}

			// If authorized, continue to the next handler
			slog.Error("Unauthorized access attempt", "userID", userID, "projectID", projectID)
			http.Error(w, "unauthorized access", http.StatusUnauthorized)
		})
	}
}

// Simulated function to check user access to a project
func userHasAccessToProject(userID int, projectName string, projectStore ProjectStore) bool {
	projects, err := projectStore.GetProjectsOwnedBy([]int{userID})
	if err != nil {
		slog.Error("Failed to get projects owned by user", "project", projectName, "user", userID, "error", err)
		return false
	}
	for _, project := range projects {
		if project.Name == projectName {
			return true
		}
	}
	return false
}
