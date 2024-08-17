package handlers

/* //Used for testing
import (
    "context"
    "math/rand"
    "net/http"
    "time"
    "encoding/json"
    "log"
    "elo-flipperz/db"
    "elo-flipperz/models"
	"go.mongodb.org/mongo-driver/bson"
)

type GameResult struct {
    Success bool `json:"success"`
    Elo     int  `json:"elo"`
	Outcome int  `json:"outcome"`
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

func FlipCoin(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    username := r.URL.Query().Get("username")

    collection := db.GetCollection("players")

    var player models.Player
    err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&player)
    if err != nil {
        http.Error(w, "Player not found", http.StatusNotFound)
        return
    }

    outcome := rand.Intn(2)

    if outcome == 1 {
        player.Elo += 10
    } else {
        player.Elo -= 10
    }

    _, err = collection.UpdateOne(
        context.TODO(),
        bson.M{"_id": player.ID},
        bson.M{"$set": bson.M{"elo": player.Elo}},
    )
    if err != nil {
        log.Fatal(err)
    }

    json.NewEncoder(w).Encode(GameResult{Success: outcome == 1, Elo: player.Elo, Outcome: outcome})
}
*/