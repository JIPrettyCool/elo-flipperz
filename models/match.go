package models

import "time"

type Match struct {
    ID        string            `json:"id"`
    Winner    Player            `json:"winner"`
    Loser     Player            `json:"loser"`
    EloBefore map[string]int    `json:"elo_before"`
    EloAfter  map[string]int    `json:"elo_after"`
    MatchTime time.Time         `json:"match_time"`
}
