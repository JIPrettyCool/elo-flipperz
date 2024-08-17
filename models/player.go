package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Username      string             `bson:"username"`
	Password	string             `bson:"password"`
    Elo       int                `bson:"elo"`
    RegisterDate time.Time          `bson:"register_date"`
    LastLoginDate time.Time          `bson:"last_login_date"`
    QueueToken  string             `bson:"queue_token"`
}