package controllers

type AppConfig struct {
	AppName   string `json:"app_name"`
	AppPort   string `json:"app_port"`
	AppURL    string `json:"app_url"`
	AppSecret string `json:"app_secret"`
}

type UserPostLogin struct {
	ID       string `bson:"_id"` // bson tag is used for MongoDB
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
