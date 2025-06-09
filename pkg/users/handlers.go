package users

/*
// Login handler
func createLoginHandler2(service UserStore) fiber.Handler {
	return func(c fiber.Ctx) error {

		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var body LoginRequest

		if err := json.Unmarshal(c.Body(), &body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Fetch user from DB

		user, err := service.GetUserByName(body.Username)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}


		// Check password
		if err := CheckPassword(user.HashedPassword, body.Password); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Generate JWT
		token, err := generateJWT(user.Name)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		// Respond with token
		return c.JSON(fiber.Map{"token": token})
	}
}
*/
/*
func generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return "", err
	}
	return token.SignedString(jwtSecret)
}

func getJWTSecret() (string, error) {
	secret := viper.GetString("jwtSecret")
	if secret == "" {
		return "", fmt.Errorf("unable to get jwt secret")
	}
	return secret, nil
}*/
