package domain

// 业务概念
type User struct {
	Id       int64
	Email    string
	Password string
	Nickname string
	Birthday string
	AboutMe  string
}
