package jwt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	userID int = 2
)

func TestJWT() error {
	userclaims := UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	signedAccessToken, err := NewAccessToken(userclaims)
	if err != nil {
		return err
	}

	refreshClaims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	signedRefreshToken, err := NewRefreshToken(refreshClaims)
	if err != nil {
		log.Fatal("error creating refresh token")
	}
	b, _ := json.MarshalIndent(signedRefreshToken, "", "\t")
	fmt.Println(string(b))

	/////

	userClaimsParsed, err := ParseAccessToken(signedAccessToken)
	if err != nil {
		return err
	}
	fmt.Println("UserID: ", userClaimsParsed.UserID)

	return nil
}
