package entities

import (
	"AwTV/modules/device/domain/entities"
	"time"
)

type Layout struct {
	Uuid       string          `json:"uuid"`
	Name       string          `json:"name"`
	Device     entities.Device `json:"device"`
	Content    string          `json:"content"`
	CreatedAt  *time.Time      `json:"created_at"`
	ModifiedAt *time.Time      `json:"modified_at"`
}
