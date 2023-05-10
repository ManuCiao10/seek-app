package controllers

import "time"

// type UserPostLogin struct {
// 	ID       string `bson:"_id"` // bson tag is used for MongoDB
// 	Email    string `form:"email" json:"email" binding:"required"`
// 	Password string `form:"password" json:"password" binding:"required"`
// }

// type UserPostSignup struct {
// 	ID        string    `bson:"_id"`
// 	Fullname  string    `form:"fullname" json:"fullname" binding:"required"`
// 	Username  string    `form:"username" json:"username" binding:"required"`
// 	Email     string    `form:"email" json:"email" binding:"required"`
// 	Password  string    `form:"password" json:"password" binding:"required"`
// 	ExpiresAt time.Time `json:"expiresAt"`
// }

// type User interface {
// 	GetUserID() string
// 	GetUserIDGoogle() string
// 	GetUserIDDiscord() string
// }

type UserPost struct {
	ID        string    `bson:"_id"`
	Fullname  string    `form:"fullname" json:"fullname" binding:"required"`
	Username  string    `form:"username" json:"username" binding:"required"`
	Email     string    `form:"email" json:"email" binding:"required"`
	Password  string    `form:"password" json:"password" binding:"required"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type UserGoogle struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	Locale     string `json:"locale"`
}

type UserDiscord struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Locale   string `json:"locale"`
	Email    string `json:"email"`
}

// Discord scope constants.
const (
	ScopeIdentify = "identify"
	ScopeBot      = "bot"
	ScopeEmail    = "email"
	ScopeGuilds   = "guilds"
)
