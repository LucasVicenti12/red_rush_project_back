package entities

import "time"

type Orientation int

const (
	Portrait  Orientation = 0
	Landscape Orientation = 1
)

func (o Orientation) IsValid() bool {
	switch o {
	case Portrait, Landscape:
		return true
	default:
		return false
	}
}

type Device struct {
	Uuid        string      `json:"uuid"`
	Name        string      `json:"name"`
	Height      int         `json:"height"`
	Width       int         `json:"width"`
	Orientation Orientation `json:"orientation"`
	CreatedAt   *time.Time  `json:"created_at"`
	ModifiedAt  *time.Time  `json:"modified_at"`
}
