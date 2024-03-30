package authentication

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func verifyToken(tokenString string, secretKey []byte) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		fmt.Printf("Error happened verifying error: %s", err)
		return err
	}

	if !token.Valid {
		fmt.Println("Invalid token:", token.Valid)
		return errors.New("invalid token")
	}

	return nil
}

func GenerateToken(username, password string) (string, time.Time, error) {
	registeredDate := time.Now().Add(time.Hour * 72)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      registeredDate.Unix(),
	})
	// generate string token
	secretKey := []byte(password)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Printf("Error happened generating string token: %s", err)
		return "", time.Time{}, err
	}

	// verify token
	err = verifyToken(tokenString, secretKey)
	if err != nil {
		fmt.Printf("Error happened calling calling verifying func: %s", err)
		return "", time.Time{}, err
	}

	return tokenString, registeredDate, nil
}
