package entity

type User struct {
	Id   int    `json:"id" bson:"id"`
	Name string `json:"userName" bson:"name"`
}
