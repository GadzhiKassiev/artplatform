package model

import "time"

type UserID = string

type Role string

const (
	RoleStudent Role = "STUDENT"
	RoleTeacher Role = "TEACHER"
)

type User struct {
	ID           UserID
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
