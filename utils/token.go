package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Name      string
	Fullname 	string
	Username  string
	Uid       string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(name string, fullname string, username string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
			Name:      	name,
			Fullname: 	fullname,
			Username:  	username,
			Uid:        uid,
			StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			},
	}

	refreshClaims := &SignedDetails{
			StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
			},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
			log.Panic(err)
			return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
			signedToken,
			&SignedDetails{},
			func(token *jwt.Token) (interface{}, error) {
					return []byte(SECRET_KEY), nil
			},
	)

	if err != nil {
			msg = err.Error()
			return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
			msg = fmt.Sprintf("the token is invalid")
			msg = err.Error()
			return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
			msg = fmt.Sprintf("token is expired")
			msg = err.Error()
			return
	}

	return claims, msg
}

