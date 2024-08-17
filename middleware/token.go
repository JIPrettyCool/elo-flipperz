package middleware

import (
    "net/http"
    "strings"
    "elo-flipperz/auth"
)

func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Where is the token amk", http.StatusUnauthorized)
            return
        }

        token = strings.TrimPrefix(token, "Bearer ")
        username, err := auth.ValidateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        r.Header.Set("username", username)
        next.ServeHTTP(w, r)
    })
}