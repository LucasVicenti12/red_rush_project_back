package use_case

import (
	"AwTV/modules/layouts/domain/entities"
	"AwTV/modules/layouts/domain/repository"
	"errors"
)

type LayoutUseCase interface {
	GetLayouts() (*[]entities.Layout, error)
	GetLayoutByUUID(uuid string) (*entities.Layout, error)
	GetLayoutByName(name string) (*[]entities.Layout, error)
	SaveLayout(layout entities.Layout) (*entities.Layout, error)
	DeleteLayout(uuid string) error
}

type LayoutUseCaseImpl struct {
	Repository repository.LayoutRepository
}

func (u LayoutUseCaseImpl) GetLayouts() (*[]entities.Layout, error) {
	return u.Repository.GetLayouts()
}

func (u LayoutUseCaseImpl) GetLayoutByUUID(uuid string) (*entities.Layout, error) {
	if uuid == "" {
		return nil, errors.New("uuid is required")
	}

	return u.Repository.GetLayoutByUUID(uuid)
}

func (u LayoutUseCaseImpl) GetLayoutByName(name string) (*[]entities.Layout, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	return u.Repository.GetLayoutByName(name)
}

func (u LayoutUseCaseImpl) SaveLayout(layout entities.Layout) (*entities.Layout, error) {
	if layout.Name == "" {
		return nil, errors.New("Name is required")
	}

	if layout.Device.Uuid == "" {
		return nil, errors.New("Device is required")
	}

	if layout.Content == "" {
		return nil, errors.New("The content is empty")
	}

	exists, err := u.Repository.ExistsLayout(layout.Uuid)

	if err != nil {
		return nil, err
	}

	if exists {
		return u.SaveLayout(layout)
	} else {
		return u.Repository.CreateLayout(layout)
	}
}

func (u LayoutUseCaseImpl) DeleteLayout(uuid string) error {
	exists, err := u.Repository.ExistsLayout(uuid)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("layout does not exist")
	}

	return u.Repository.DeleteLayout(uuid)
}
