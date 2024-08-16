package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    "log"

    "go.mongodb.org/mongo-driver/bson/primitive"

    "elo-flipperz/db"
    "elo-flipperz/models"
)

type CreatePlayerRequest struct {
    Username string `json:"name"`
	Password string `json:"password"`
}

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request for CreatePlayer")
    w.Header().Set("Content-Type", "application/json")

    var req CreatePlayerRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    player := models.Player{
        ID:           primitive.NewObjectID(),
        Username:     req.Username,
        Password:     req.Password,
        Elo:          1000,
        RegisterDate: time.Now(),
    }
    collection := db.GetCollection("players")
    _, err = collection.InsertOne(context.TODO(), player)
    if err != nil {
        log.Fatal(err)
    }

    json.NewEncoder(w).Encode(player)
}