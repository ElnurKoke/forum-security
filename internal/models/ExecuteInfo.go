package models

type Info struct {
	User
	Post
	Comment     []Comment
	AllCategory []Category
}

type InfoPosts struct {
	User
	Posts    []Post
	Category []Category
}

type InfoMsg struct {
	User
	Notifications []Message
	Actions       []Message
}

type InfoSign struct {
	Error          string
	Username       string
	Password       string
	RepeatPassword string
	Email          string
}
