package controllers

type AppConfig struct {
	AppName string `json:"app_name"`
	AppPort string `json:"app_port"`
	AppURL  string `json:"app_url"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
