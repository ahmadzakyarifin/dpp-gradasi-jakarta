package main

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := "ubah-ini-dengan-secret-panjang-minimal-32-karakter"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"role_name":    "super_admin",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	fmt.Println(tokenString)
}
