package repository

import (
	"AwTV/modules/device/domain/entities"
	"database/sql"
	"errors"
	"time"
)

type DeviceRepositoryImpl struct {
	db *sql.DB
}

func NewDeviceRepository(conn *sql.DB) *DeviceRepositoryImpl {
	return &DeviceRepositoryImpl{
		db: conn,
	}
}

func (dr *DeviceRepositoryImpl) GetDeviceByName(name string) (*[]entities.Device, error) {
	devices := make([]entities.Device, 0)

	rows, err := dr.db.Query(
		"SELECT d.uuid, d.name, d.height, d.width, d.orientation, d.created_at, d.modified_at FROM devices d WHERE name like $1",
		name,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &devices, err
	}

	defer func() {
		_ = rows.Close()
	}()

	if rows.Next() {
		var device entities.Device

		err := rows.Scan(
			&device.Uuid,
			&device.Name,
			&device.Height,
			&device.Width,
			&device.Orientation,
			&device.CreatedAt,
			&device.ModifiedAt,
		)

		if err != nil {
			return &devices, err
		}

		devices = append(devices, device)
	}

	return &devices, nil
}

func (dr *DeviceRepositoryImpl) GetDeviceByUUID(uuid string) (*entities.Device, error) {
	rows, err := dr.db.Query(
		"SELECT d.uuid, d.name, d.height, d.width, d.orientation, d.created_at, d.modified_at FROM devices d WHERE d.uuid = $1",
		uuid,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var device entities.Device

	if rows.Next() {
		err := rows.Scan(
			&device.Uuid,
			&device.Name,
			&device.Height,
			&device.Width,
			&device.Orientation,
			&device.CreatedAt,
			&device.ModifiedAt,
		)

		if err != nil {
			return nil, err
		}
	}

	return &device, nil
}

func (dr *DeviceRepositoryImpl) GetAllDevices() (*[]entities.Device, error) {
	devices := make([]entities.Device, 0)

	rows, err := dr.db.Query("SELECT d.uuid, d.name, d.height, d.width, d.orientation, d.created_at, d.modified_at FROM devices d")

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &devices, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var device entities.Device

		var created sql.NullString
		var modified sql.NullString

		err := rows.Scan(
			&device.Uuid,
			&device.Name,
			&device.Height,
			&device.Width,
			&device.Orientation,
			&created,
			&modified,
		)

		if created.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", created.String)

			device.CreatedAt = &t
		}

		if modified.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", modified.String)

			device.ModifiedAt = &t
		}

		if err != nil {
			return &devices, err
		}

		devices = append(devices, device)
	}

	return &devices, nil
}

func (dr *DeviceRepositoryImpl) ExistsDevice(uuid string) (bool, error) {
	rows, err := dr.db.Query(
		"SELECT d.uuid FROM devices d WHERE d.uuid = ?",
		uuid,
	)

	if err != nil {
		return false, err
	}

	defer func() {
		_ = rows.Close()
	}()

	r := ""

	if rows.Next() {
		err := rows.Scan(&r)

		if err != nil {
			return false, err
		}
	}

	return r != "", nil
}

func (dr *DeviceRepositoryImpl) CreateDevice(entity *entities.Device) error {
	smt, err := dr.db.Prepare("INSERT INTO devices (uuid, name, height, width, orientation, created_at, modified_at) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}

	time.Local = time.UTC

	currTime := time.Now()

	if entity.CreatedAt == nil {
		entity.CreatedAt = &currTime
	}

	entity.ModifiedAt = &currTime

	_, err = smt.Exec(
		entity.Uuid,
		entity.Name,
		entity.Height,
		entity.Width,
		entity.Orientation,
		entity.CreatedAt,
		entity.ModifiedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (dr *DeviceRepositoryImpl) EditDevice(entity *entities.Device) error {
	smt, err := dr.db.Prepare("UPDATE devices SET name = ?, height = ?, width = ?, orientation = ?, modified_at = ? WHERE uuid = ?")

	if err != nil {
		return err
	}

	time.Local = time.UTC
	currTime := time.Now()

	entity.ModifiedAt = &currTime

	_, err = smt.Exec(
		entity.Name,
		entity.Height,
		entity.Width,
		entity.Orientation,
		entity.ModifiedAt,
		entity.Uuid,
	)

	if err != nil {
		return err
	}

	return nil
}

func (dr *DeviceRepositoryImpl) DeleteDevice(uuid string) error {
	stm, err := dr.db.Prepare("DELETE FROM devices d WHERE d.uuid = ?")

	if err != nil {
		return err
	}

	_, err = stm.Exec(uuid)

	if err != nil {
		return err
	}

	return nil
}
