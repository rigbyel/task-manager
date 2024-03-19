package model

type User struct {
	Id      int64
	Name    string
	Balance int
}

type Quest struct {
	Id   int64
	Name string
	Cost int
}
