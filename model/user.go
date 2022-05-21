package model

import "context"

type User struct {
	Id       int64
	Name     string
	Email    string
	Password string
}

type AccountRepository interface {
	UserCreate(ctx context.Context, Name string, Email string, Password string) (*User, error)
	UserGetByID(ctx context.Context, id string) (*User, error)
	UserGetByEmail(ctx context.Context, email string) (*User, error)
}
