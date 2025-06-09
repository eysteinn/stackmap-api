package auth

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
