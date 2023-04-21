package controllers

import "time"

type AppConfig struct {
	AppName   string `json:"app_name"`
	AppPort   string `json:"app_port"`
	AppURL    string `json:"app_url"`
	AppSecret string `json:"app_secret"`
}

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Rank        string    `json:"rank"`
	Reviews     string    `json:"reviews"`
	Descritpion string    `json:"description"`
	Location    string    `json:"location"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	Lastseen    time.Time `json:"lastseen"`
}
