package models

type User struct {
	Password string `json:"password", gorm:"password"`
	Username string `json:"username", gorm:"username"`
}
