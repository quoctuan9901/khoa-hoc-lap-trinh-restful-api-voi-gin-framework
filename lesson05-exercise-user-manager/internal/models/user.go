package models

type User struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Level    int    `json:"level"`
}
