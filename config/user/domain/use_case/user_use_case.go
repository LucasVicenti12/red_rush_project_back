package use_case

import (
	"AwTV/config/user/domain/entities"
	"AwTV/config/user/domain/repository"
	"errors"
	uuid2 "github.com/google/uuid"
)

type UserUseCase interface {
	GetUserByUUID(uuid string) (*entities.User, error)
	GetUserByNickname(nickname string) (*entities.User, error)
	RegisterUser(user entities.User) error
	ListUsers() (*[]entities.User, error)
}

type UserUseCaseImpl struct {
	Repository repository.UserRepository
}

func (u UserUseCaseImpl) GetUserByUUID(uuid string) (*entities.User, error) {
	return u.Repository.GetUserByUUID(uuid)
}

func (u UserUseCaseImpl) GetUserByNickname(nickname string) (*entities.User, error) {
	return u.Repository.GetUserByNickname(nickname)
}

func (u UserUseCaseImpl) RegisterUser(user entities.User) error {
	us, err := u.Repository.GetUserByEmail(user.Email)

	if err != nil {
		return err
	}

	if us != nil {
		return errors.New("user with this e-mail already exists")
	}

	us, err = u.Repository.GetUserByNickname(user.Nickname)

	if err != nil {
		return err
	}

	if us != nil {
		return errors.New("user with this nickname already exists")
	}

	uuid, _ := uuid2.NewUUID()

	user.Uuid = uuid.String()

	return u.Repository.RegisterUser(user)
}

func (u UserUseCaseImpl) ListUsers() (*[]entities.User, error) {
	return u.Repository.ListUsers()
}
