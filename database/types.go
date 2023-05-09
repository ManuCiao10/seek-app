package database

import "time"

var user UserDatabase

const (
	dbName     = "GODB"
	collectionName = "account"
)

type UserDatabase struct {
	ID        string    `bson:"_id"`
	Fullname  string    `form:"fullname" json:"fullname" binding:"required"`
	Username  string    `form:"username" json:"username" binding:"required"`
	Email     string    `form:"email" json:"email" binding:"required"`
	Password  string    `form:"password" json:"password" binding:"required"`
	Image     string    `json:"image"`
	Locale    string    `json:"locale"`
	Ranking   int       `json:"ranking"`
	SessionID string    `json:"sessionID"`
	ExpiresAt time.Time `json:"expiresAt"`
}
