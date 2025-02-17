package models

import "time"

type User struct {
	Id           int
	Name         string
	Email        string
	PasswordHash string
	PhoneNumber  *int
	Birthday     *time.Time
}


type Roles struct {
	Id           int
	Name         string
	Email        string
	PasswordHash string
	PhoneNumber  *int
	Birthday     *time.Time
}
