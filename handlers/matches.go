package handlers

import (
    "net/http"
    "strings"
    "context"
    "encoding/json"
    "go.mongodb.org/mongo-driver/bson"
    "elo-flipperz/db"
    "elo-flipperz/models"
)

func HandleMatchResult(w http.ResponseWriter, r *http.Request) {
    if strings.HasPrefix(r.URL.Path, "/matches/") {
        matchID := strings.TrimPrefix(r.URL.Path, "/matches/")
        if matchID == "" {
            http.Error(w, "Match ID is required", http.StatusBadRequest)
            return
        }
        r.URL.RawQuery = "id=" + matchID
        GetMatchResult(w, r)
    } else {
        http.NotFound(w, r)
    }
}

func GetMatchResult(w http.ResponseWriter, r *http.Request) {
    matchID := r.URL.Path[len("/matches/"):]

    if matchID == "" {
        http.Error(w, "where is Match ID biç", http.StatusBadRequest)
        return
    }

    collection := db.GetCollection("matches")
    var match models.Match

    err := collection.FindOne(context.TODO(), bson.M{"id": matchID}).Decode(&match)
    if err != nil {
        http.Error(w, "There is no match with that ID", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(match)
}