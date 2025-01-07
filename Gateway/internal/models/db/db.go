package db

type User struct {
	ID       uint64
	Login    string
	Password string
	Role     string
}
