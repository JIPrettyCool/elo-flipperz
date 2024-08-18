package handlers

import (
    "context"
    "fmt"
    "encoding/json"
    "math/rand"
    "net/http"
    "strings"
    "sync"
    "time"

    "elo-flipperz/auth"
    "elo-flipperz/db"
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
        matchID := storeMatchResult(result)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":    "match completed",
            "match_id":  matchID,
            "result":    result,
        })
    } else {
        json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
    }
}

func storeMatchResult(result map[string]interface{}) string {
    winner, ok1 := result["winner"].(models.Player)
    loser, ok2 := result["loser"].(models.Player)

    fmt.Println(winner, ok1, loser, ok2) // TODO: remove this add debug messages instead
    
    matchID := generateMatchID(winner.Username, loser.Username)
    
    match := models.Match{
        ID:        matchID,
        Winner:    winner, // TODO: change it to Username before full build
        Loser:     loser, // TODO: change it to Username before full build
        EloBefore: map[string]int{
            winner.Username: result["elo_before_winner"].(int),
            loser.Username:  result["elo_before_loser"].(int),
        },
        EloAfter: map[string]int{
            winner.Username: winner.Elo + 10,
            loser.Username:  loser.Elo - 10,
        },
        MatchTime: time.Now(),
    }

    collection := db.GetCollection("matches")
    _, err := collection.InsertOne(context.TODO(), match)
    if err != nil {
        panic("Failed to insert match: " + err.Error())
    }

    return matchID
}

func generateMatchID(winnerUsername, loserUsername string) string {
    timestamp := time.Now().Format("20060102150405")
    return winnerUsername + "_" + loserUsername + "_" + timestamp
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
    
    eloBeforeWinner := winner.Elo
    eloBeforeLoser := loser.Elo

    updateElo(winner, loser)

    return map[string]interface{}{
        "status":             "match completed",
        "winner":             winner,
        "loser":              loser,
        "elo_before_winner":  eloBeforeWinner,
        "elo_before_loser":   eloBeforeLoser,
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