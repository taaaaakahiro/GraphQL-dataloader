package model

type Message struct {
	ID string `json:"id"`
	// User    *User  `json:"user"`
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}
