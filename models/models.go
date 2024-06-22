package models

type User struct {
	// gorm.Model

	User_Id       string `json:"user_id"`
	User_name     string `json:"user_name"`
	User_Email    string `json:"user_email"`
	User_password string `json:"-"`
}
