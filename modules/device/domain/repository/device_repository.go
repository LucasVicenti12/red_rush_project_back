package repository

import (
	"AwTV/modules/device/domain/entities"
)

type DeviceRepository interface {
	GetDeviceByName(name string) (*[]entities.Device, error)
	GetDeviceByUUID(uuid string) (*entities.Device, error)
	GetAllDevices() (*[]entities.Device, error)
	ExistsDevice(uuid string) (bool, error)
	CreateDevice(entity *entities.Device) error
	EditDevice(entity *entities.Device) error
	DeleteDevice(uuid string) error
}
