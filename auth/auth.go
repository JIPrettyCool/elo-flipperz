package auth

import (
    "github.com/dgrijalva/jwt-go"
    "time"
)

var jwtSecret = []byte("your_secret_key") // Use a secure key

func GenerateToken(username string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "expire":      time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (string, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || !token.Valid {
        return "", err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return "", err
    }

    username, ok := claims["username"].(string)
    if !ok {
        return "", err
    }

    return username, nil
}