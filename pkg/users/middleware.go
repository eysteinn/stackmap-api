package users

import (
	"net/http"
)

type jwtMiddleware struct {
	Userstore *UserStore
}

func (mw jwtMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
func CreateJWTMiddleware(userstore *UserStore) *jwtMiddleware {
	return &jwtMiddleware{
		Userstore: userstore,
	}
}

/*
// Middleware to validate JWT and extract userID
func jwtMiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, `{"error": "Missing or invalid Authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Extract userID from token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, ok := claims["userID"].(string)
			if !ok {
				http.Error(w, `{"error": "Invalid token payload"}`, http.StatusUnauthorized)
				return
			}

			// Add userID to request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
	})
}

func test() {
	mw := *CreateJWTMiddleware(nil)
	a := http.Handler(mw)

	alice.New(a)
}
*/
