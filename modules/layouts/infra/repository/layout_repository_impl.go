package repository

import (
	"AwTV/modules/device/domain/repository"
	repository2 "AwTV/modules/device/infra/repository"
	"AwTV/modules/layouts/domain/entities"
	"database/sql"
	"time"
)

type LayoutRepositoryImpl struct {
	db               *sql.DB
	deviceRepository repository.DeviceRepository
}

func NewLayoutRepositoryImpl(db *sql.DB) *LayoutRepositoryImpl {
	return &LayoutRepositoryImpl{
		db:               db,
		deviceRepository: repository2.NewDeviceRepository(db),
	}
}

func (lr *LayoutRepositoryImpl) GetLayouts() (*[]entities.Layout, error) {
	layouts := make([]entities.Layout, 0)

	rows, err := lr.db.Query("SELECT l.uuid, l.name, l.device_uuid, l.content, l.created_at, l.modified_at FROM layouts l")

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var layout entities.Layout
		var deviceUUID string

		var created sql.NullString
		var modified sql.NullString

		err := rows.Scan(
			&layout.Uuid,
			&layout.Name,
			&deviceUUID,
			&created,
			&modified,
		)

		if created.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", created.String)

			layout.CreatedAt = &t
		}

		if modified.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", modified.String)

			layout.ModifiedAt = &t
		}

		if err != nil {
			return nil, err
		}

		device, err := lr.deviceRepository.GetDeviceByUUID(deviceUUID)

		layout.Device = *device

		if err != nil {
			return nil, err
		}

		layouts = append(layouts, layout)
	}

	return &layouts, nil
}

func (lr *LayoutRepositoryImpl) GetLayoutByUUID(uuid string) (*entities.Layout, error) {
	rows, err := lr.db.Query(
		"SELECT l.uuid, l.name, l.device_uuid, l.content, l.created_at, l.modified_at FROM layouts l WHERE l.uuid = $1",
		uuid,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var layout entities.Layout
	if rows.Next() {
		var deviceUUID string

		var created sql.NullString
		var modified sql.NullString

		err := rows.Scan(
			&layout.Uuid,
			&layout.Name,
			&deviceUUID,
			&created,
			&modified,
		)

		if err != nil {
			return nil, err
		}

		if created.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", created.String)

			layout.CreatedAt = &t
		}

		if modified.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", modified.String)

			layout.ModifiedAt = &t
		}

		device, err := lr.deviceRepository.GetDeviceByUUID(deviceUUID)

		layout.Device = *device

		if err != nil {
			return nil, err
		}
	}

	return &layout, nil
}

func (lr *LayoutRepositoryImpl) GetLayoutByName(name string) (*[]entities.Layout, error) {
	layouts := make([]entities.Layout, 0)

	rows, err := lr.db.Query(
		"SELECT l.uuid, l.name, l.device_uuid, l.content, l.created_at, l.modified_at FROM layouts l WHERE l.name like $1",
		name,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var layout entities.Layout
		var deviceUUID string

		var created sql.NullString
		var modified sql.NullString

		err := rows.Scan(
			&layout.Uuid,
			&layout.Name,
			&deviceUUID,
			&created,
			&modified,
		)

		if err != nil {
			return nil, err
		}

		if created.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", created.String)

			layout.CreatedAt = &t
		}

		if modified.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", modified.String)

			layout.ModifiedAt = &t
		}

		device, err := lr.deviceRepository.GetDeviceByUUID(deviceUUID)

		layout.Device = *device

		if err != nil {
			return nil, err
		}

		layouts = append(layouts, layout)
	}

	return &layouts, nil
}

func (lr *LayoutRepositoryImpl) CreateLayout(layout entities.Layout) (*entities.Layout, error) {
	smt, err := lr.db.Prepare("INSERT INTO layouts (uuid, name, device_uuid, content, created_at, modified_at) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = smt.Close()
	}()

	time.Local = time.UTC

	currTime := time.Now()

	if layout.CreatedAt == nil {
		layout.CreatedAt = &currTime
	}

	layout.ModifiedAt = &currTime

	_, err = smt.Exec(
		layout.Uuid,
		layout.Name,
		layout.Device.Uuid,
		layout.Content,
		layout.CreatedAt,
		layout.ModifiedAt,
	)

	if err != nil {
		return nil, err
	}

	return &layout, nil
}

func (lr *LayoutRepositoryImpl) EditLayout(layout entities.Layout) (*entities.Layout, error) {
	smt, err := lr.db.Prepare("UPDATE layouts SET name = ?, content = ?, modified_at = ? WHERE uuid = ?")

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = smt.Close()
	}()

	time.Local = time.UTC

	currTime := time.Now()

	layout.ModifiedAt = &currTime

	_, err = smt.Exec(
		layout.Name,
		layout.Content,
		layout.ModifiedAt,
		layout.Uuid,
	)

	if err != nil {
		return nil, err
	}

	return &layout, nil
}

func (lr *LayoutRepositoryImpl) ExistsLayout(uuid string) (bool, error) {
	rows, err := lr.db.Query(
		"SELECT l.uuid FROM layouts l WHERE l.uuid = ?",
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

func (lr *LayoutRepositoryImpl) DeleteLayout(uuid string) error {
	stm, err := lr.db.Prepare("DELETE FROM layouts l WHERE l.uuid = ?")

	if err != nil {
		return err
	}

	_, err = stm.Exec(uuid)

	if err != nil {
		return err
	}

	return nil
}
