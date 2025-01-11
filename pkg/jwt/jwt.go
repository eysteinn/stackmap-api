package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/utils"
)

var access_token_expiration time.Duration = time.Minute * 1000 // should be 15 minutes

type UserClaims struct {
	UserID int `json:"user_id"`
	//First string `json:"first"`
	//Last  string `json:"last"`
	jwt.StandardClaims
}

func InitUserClaims(userID int) UserClaims {
	userclaims := UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(access_token_expiration).Unix(),
		},
	}
	return userclaims
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := GetJWTSecret()
	if err != nil {
		return "", err
	}
	return accessToken.SignedString(secret)
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := GetJWTSecret()
	if err != nil {
		return "", err
	}
	return refreshToken.SignedString(secret)
}

func ParseAccessToken(accessToken string) (*UserClaims, error) {

	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret()
	})
	if err != nil {
		return nil, err
	}
	return parsedAccessToken.Claims.(*UserClaims), nil
}

func ParseRefreshToken(refreshToken string) (*jwt.StandardClaims, error) {

	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret()
	})
	if err != nil {
		return nil, err
	}

	return parsedRefreshToken.Claims.(*jwt.StandardClaims), nil
}

/////////////

//https://pascalallen.medium.com/jwt-authentication-with-go-242215a9b4f8
/*
type UserData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type JWTClaim struct {
	Data UserData `json:"data"`
	jwt.StandardClaims
}
*/

func GetJWTSecret() ([]byte, error) {
	secret := viper.GetString(utils.JWT_SECRET_KEY)
	if secret == "" {
		return nil, fmt.Errorf("unable to get jwt secret")
	}
	return []byte(secret), nil
}

func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret, err := GetJWTSecret()
	if err != nil {
		return "", err
	}
	return token.SignedString(jwtSecret)
}

/*
	func ValidateJWT2(tokenString string) (*jwt.Token, error) {
		secretKey, err := GetJWTSecret()
		if err != nil {
			return nil, err
		}

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{},
			func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

		return token, err
	}
*/
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKey, err := GetJWTSecret()
		if err != nil {
			return "", nil
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check token validity
	/*if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid")
		fmt.Println("Claims:", claims)
	} else {
		return nil, fmt.Errorf("invalid token")
	}*/

	return token, nil
}
