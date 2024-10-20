package repository

import "AwTV/config/user/domain/entities"

type UserRepository interface {
	GetUserByUUID(uuid string) (*entities.User, error)
	GetUserByNickname(nickname string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	RegisterUser(user entities.User) error
	ListUsers() (*[]entities.User, error)
}
