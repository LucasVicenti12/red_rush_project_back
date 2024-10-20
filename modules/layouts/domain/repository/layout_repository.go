package repository

import "AwTV/modules/layouts/domain/entities"

type LayoutRepository interface {
	GetLayouts() (*[]entities.Layout, error)
	GetLayoutByUUID(uuid string) (*entities.Layout, error)
	GetLayoutByName(name string) (*[]entities.Layout, error)
	CreateLayout(layout entities.Layout) (*entities.Layout, error)
	EditLayout(layout entities.Layout) (*entities.Layout, error)
	ExistsLayout(uuid string) (bool, error)
	DeleteLayout(uuid string) error
}
