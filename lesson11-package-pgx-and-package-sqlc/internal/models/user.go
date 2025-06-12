package models

type User struct {
	Id    int    `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
