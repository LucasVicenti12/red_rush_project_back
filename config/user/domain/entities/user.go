package entities

import "time"

type UserType int

const (
	Admin      UserType = 0
	DeviceUser UserType = 1
)

type User struct {
	Uuid       string     `json:"uuid"`
	Name       string     `json:"name"`
	Nickname   string     `json:"nickname"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	UserType   UserType   `json:"user_type"`
	CreatedAt  *time.Time `json:"created_at"`
	ModifiedAt *time.Time `json:"modified_at"`
}
