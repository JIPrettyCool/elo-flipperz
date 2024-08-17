package handlers

import (
    "net/http"
    "sync"
    "math/rand"
    "context"
    "encoding/json"
    "strings"

    "elo-flipperz/db"
    "elo-flipperz/auth"
    "elo-flipperz/models"
    "go.mongodb.org/mongo-driver/bson"
)

var (
    queue     []models.Player
    queueLock sync.Mutex
)

func QueuePlayer(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    token := r.Header.Get("Authorization")
    if token == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    username, err := auth.ValidateToken(strings.TrimPrefix(token, "Bearer "))
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    collection := db.GetCollection("players")
    var player models.Player
    err = collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&player)
    if err != nil {
        http.Error(w, "Player not found", http.StatusNotFound)
        return
    }

    queueLock.Lock()
    queue = append(queue, player)
    queueLock.Unlock()

    if len(queue) >= 2 {
        result := StartMatch()
        json.NewEncoder(w).Encode(result)
    } else {
        json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
    }
}

func StartMatch() map[string]interface{} {
    queueLock.Lock()
    defer queueLock.Unlock()
    player1 := queue[0]
    player2 := queue[1]
    queue = queue[2:]

    if player1.Elo == 0 || player2.Elo == 0 {
        return map[string]interface{}{"status": "Bro you can't play with 0 Elo"}
    }
    if player1 == player2 {
        return map[string]interface{}{"status": "same player"}
    }

    winner, loser := determineWinner(player1, player2)
    updateElo(winner, loser)

    return map[string]interface{}{
        "status": "match completed",
        "winner": map[string]interface{}{
            "username": winner.Username,
            "elo":      winner.Elo,
        },
        "loser": map[string]interface{}{
            "username": loser.Username,
            "elo":      loser.Elo,
        },
        "result": 1,
    }
}

func determineWinner(player1, player2 models.Player) (models.Player, models.Player) {
    if rand.Intn(2) == 0 {
        return player1, player2
    }
    return player2, player1
}

func updateElo(winner, loser models.Player) {
    collection := db.GetCollection("players")

    winner.Elo += 10
    loser.Elo -= 10

    collection.UpdateOne(context.TODO(), bson.M{"_id": winner.ID}, bson.M{"$set": bson.M{"elo": winner.Elo}})
    collection.UpdateOne(context.TODO(), bson.M{"_id": loser.ID}, bson.M{"$set": bson.M{"elo": loser.Elo}})
}