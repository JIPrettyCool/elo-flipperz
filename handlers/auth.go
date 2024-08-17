package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "elo-flipperz/db"
    "elo-flipperz/models"
    "elo-flipperz/auth"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type RegistrationData struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginData struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var data LoginData
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    collection := db.GetCollection("players")
    var player models.Player
    err := collection.FindOne(context.TODO(), bson.M{"username": data.Username, "password": data.Password}).Decode(&player)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    _, err = collection.UpdateOne(
        context.TODO(),
        bson.M{"_id": player.ID},
        bson.M{"$set": bson.M{"last_login_date": time.Now()}},
    )
    if err != nil {
        http.Error(w, "Error updating last login date", http.StatusInternalServerError)
        return
    }

    token, err := auth.GenerateToken(data.Username)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    _, err = collection.UpdateOne(
        context.TODO(),
        bson.M{"_id": player.ID},
        bson.M{"$set": bson.M{"queue_token": token}},
    )
    if err != nil {
        http.Error(w, "Error updating token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func Logout(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func Register(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var data RegistrationData
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    collection := db.GetCollection("players")

    var existingPlayer models.Player
    err := collection.FindOne(context.TODO(), bson.M{"username": data.Username}).Decode(&existingPlayer)
    if err == nil {
        http.Error(w, "Username already taken", http.StatusBadRequest)
        return
    } else if err != mongo.ErrNoDocuments {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    player := models.Player{
        ID:             primitive.NewObjectID(),
        Username:       data.Username,
        Password:       data.Password,
        Elo:            1000,
        RegisterDate:   time.Now(),
        LastLoginDate:  time.Time{},
        QueueToken:     "",
    }

    _, err = collection.InsertOne(context.TODO(), player)
    if err != nil {
        http.Error(w, "Could not create player", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(player)
}