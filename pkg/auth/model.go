package auth

// ProjectStore represents the project service interface
type RefreshTokenStore interface {
	ValidateRefreshToken(token string) (int, error)
	GenerateRefreshToken(userID int) (string, error)
	InvalidateRefreshToken(token string) error
}
