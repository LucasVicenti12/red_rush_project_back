package use_case

import (
	"AwTV/modules/device/domain/entities"
	"AwTV/modules/device/domain/repository"
	"errors"
	uuid2 "github.com/google/uuid"
)

type DeviceUseCase interface {
	GetDeviceByName(name string) (*[]entities.Device, error)
	GetDeviceByUUID(uuid string) (*entities.Device, error)
	GetAllDevices() (*[]entities.Device, error)
	SaveDevice(entity *entities.Device) error
	DeleteDevice(uuid string) error
}

type DeviceUseCaseImpl struct {
	Repository repository.DeviceRepository
}

func (dvuc *DeviceUseCaseImpl) GetDeviceByName(name string) (*[]entities.Device, error) {
	if name == "" {
		return nil, errors.New("device name is required")
	}

	return dvuc.Repository.GetDeviceByName(name)
}

func (dvuc *DeviceUseCaseImpl) GetDeviceByUUID(uuid string) (*entities.Device, error) {
	if uuid == "" {
		return nil, errors.New("device uuid is required")
	}

	return dvuc.Repository.GetDeviceByUUID(uuid)
}

func (dvuc *DeviceUseCaseImpl) GetAllDevices() (*[]entities.Device, error) {
	return dvuc.Repository.GetAllDevices()
}

func (dvcu *DeviceUseCaseImpl) SaveDevice(entity *entities.Device) error {
	if entity.Name == "" {
		return errors.New("device name is required")
	}

	if entity.Height == 0 {
		return errors.New("device height is required")
	}

	if entity.Width == 0 {
		return errors.New("device width is required")
	}

	if !entity.Orientation.IsValid() {
		return errors.New("device orientation is invalid")
	}

	exists, err := dvcu.Repository.ExistsDevice(entity.Uuid)

	if err != nil {
		return errors.New("An unexpected error has occurred")
	}

	if exists {
		return dvcu.Repository.EditDevice(entity)
	} else {
		uuid, _ := uuid2.NewUUID()

		entity.Uuid = uuid.String()

		return dvcu.Repository.CreateDevice(entity)
	}
}

func (dvuc *DeviceUseCaseImpl) DeleteDevice(uuid string) error {
	err := dvuc.Repository.DeleteDevice(uuid)

	return err
}
