package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Editor string `json:"editor"`
	jwt.StandardClaims
}

func EncodeJwtToken(claims *UserClaims, method jwt.SigningMethod, key string) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	tokenStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func DecodeJwtToken(tokenStr string, method jwt.SigningMethod, key string) (*jwt.Token, *UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != method.Alg() {
			return nil, fmt.Errorf("unexpedted jwt alg: got=%q, want=%q", token.Method.Alg(), method.Alg())
		}
		return []byte(key), nil
	})
	return token, claims, err
}

func main() {
	tokenStr, err := EncodeJwtToken(&UserClaims{
		Name:   "buzztaiki",
		Email:  "buzz.taiki@gmail.com",
		Editor: "GNU Emacs",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Second).Unix(),
			Issuer:    "test",
		},
	}, jwt.SigningMethodHS256, "moo")
	if err != nil {
		log.Fatalf("failed to encode token: %v", err)
	}
	log.Println(tokenStr)

	token, claims, err := DecodeJwtToken(tokenStr, jwt.SigningMethodHS256, "moo")
	if !token.Valid {
		log.Print("invalid token")
	}
	if err != nil {
		log.Fatalf("failed to decode token: %v", err)
	}
	log.Println(claims)
}
