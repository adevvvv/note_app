package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
