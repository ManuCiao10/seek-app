package controllers

type AppConfig struct {
	AppPort   string `json:"app_port"`
	AppSecret string `json:"app_secret"`
}

type UserPostLogin struct {
	ID       string `bson:"_id"` // bson tag is used for MongoDB
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserPostSignup struct {
	ID       string `bson:"_id"` // bson tag is used for MongoDB
	Fullname string `form:"fullname" json:"fullname" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
