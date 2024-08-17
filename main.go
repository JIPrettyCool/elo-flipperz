package main

import (
    "log"
    "net/http"
    "time"

    "elo-flipperz/db"
    "elo-flipperz/handlers"
    "elo-flipperz/middleware"
)

func main() {
    db.Connect()

    mux := http.NewServeMux()

    mux.Handle("/queue", middleware.ValidateToken(http.HandlerFunc(handlers.QueuePlayer)))
    mux.Handle("/register", middleware.RateLimit(http.HandlerFunc(handlers.Register)))
    mux.Handle("/login", middleware.RateLimit(http.HandlerFunc(handlers.Login)))
    mux.Handle("/leaderboard", middleware.RateLimit(http.HandlerFunc(handlers.Leaderboard)))


     corsHandler := middleware.CORS(mux)

     server := &http.Server{
         Addr:         ":8080",
         Handler:      corsHandler,
         ReadTimeout:  10 * time.Second,
         WriteTimeout: 10 * time.Second,
         IdleTimeout:  60 * time.Second,
     }
 
     log.Println("Starting server on :8080")
     log.Fatal(server.ListenAndServe())
 }