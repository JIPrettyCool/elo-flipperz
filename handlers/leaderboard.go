package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "sort"

    "elo-flipperz/db"
    "elo-flipperz/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type LeaderboardEntry struct {
    Username string `json:"username"`
    Elo      int    `json:"elo"`
}

func Leaderboard(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    collection := db.GetCollection("players")

    var leaderboard []LeaderboardEntry

    opts := options.Find().SetSort(bson.D{{Key: "elo", Value: -1}}).SetLimit(10)
    cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
    if err != nil {
        http.Error(w, "Could not retrieve leaderboard", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var player models.Player
        if err := cursor.Decode(&player); err != nil {
            http.Error(w, "Error decoding player data", http.StatusInternalServerError)
            return
        }
        leaderboard = append(leaderboard, LeaderboardEntry{
            Username: player.Username,
            Elo:      player.Elo,
        })
    }

    sort.SliceStable(leaderboard, func(i, j int) bool {
        return leaderboard[i].Elo > leaderboard[j].Elo
    })

    json.NewEncoder(w).Encode(leaderboard)
}