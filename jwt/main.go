package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JWT() (string, error) {
	mySigningKey := []byte("your-256-bit-secret")
	type MyCustomClaims struct {
		Name string `json:"name"`
		jwt.RegisteredClaims
	}
	claims := MyCustomClaims{
		Name: "John Doe",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  "1234567890",
			IssuedAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("%+v\n", *token)
	return token.SignedString(mySigningKey)
}

func main() {
	ss, _ := JWT()
	fmt.Printf("%v\n", ss)
}
