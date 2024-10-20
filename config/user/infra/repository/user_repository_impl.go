package repository

import (
	"AwTV/config/user/domain/entities"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{
		db: conn,
	}
}

func (ur UserRepository) GetUserByUUID(uuid string) (*entities.User, error) {
	query := `
	SELECT u.uuid, u.name, u.nickname, u.email, u.user_type, u.created_at, u.modified_at
	FROM users u
	WHERE u.uuid = ?
	`

	smtUser, err := ur.db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := smtUser.Close()
		if err != nil {
			return
		}
	}()

	var user entities.User

	err = smtUser.QueryRow(uuid).Scan(
		&user.Uuid,
		&user.Name,
		&user.Nickname,
		&user.Email,
		&user.UserType,
		&user.CreatedAt,
		&user.ModifiedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (ur UserRepository) GetUserByNickname(nickname string) (*entities.User, error) {
	query := `
	SELECT u.uuid, u.name, u.nickname, u.email, u.user_type, u.created_at, u.modified_at, u.password
	FROM users u
	WHERE u.nickname = ?
	`

	smtUser, err := ur.db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := smtUser.Close()
		if err != nil {
			return
		}
	}()

	var user entities.User

	var created sql.NullString
	var modified sql.NullString

	err = smtUser.QueryRow(nickname).Scan(
		&user.Uuid,
		&user.Name,
		&user.Nickname,
		&user.Email,
		&user.UserType,
		&created,
		&modified,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if created.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", created.String)

		user.CreatedAt = &t
	}

	if modified.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", modified.String)

		user.ModifiedAt = &t
	}

	return &user, nil
}

func (ur UserRepository) GetUserByEmail(email string) (*entities.User, error) {
	query := `
	SELECT u.uuid, u.name, u.nickname, u.email, u.user_type, u.created_at, u.modified_at
	FROM users u
	WHERE u.email = ?
	`

	smtUser, err := ur.db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := smtUser.Close()
		if err != nil {
			return
		}
	}()

	var user entities.User

	err = smtUser.QueryRow(email).Scan(
		&user.Uuid,
		&user.Name,
		&user.Nickname,
		&user.Email,
		&user.UserType,
		&user.CreatedAt,
		&user.ModifiedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (ur UserRepository) RegisterUser(user entities.User) error {
	query := `
	INSERT INTO users (uuid, name, nickname, email, password, user_type, created_at, modified_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	smtUser, err := ur.db.Prepare(query)

	if err != nil {
		return err
	}

	defer func() {
		err := smtUser.Close()
		if err != nil {
			return
		}
	}()

	time.Local = time.UTC

	currTime := time.Now()

	user.CreatedAt = &currTime
	user.ModifiedAt = &currTime

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	_, err = smtUser.Exec(
		user.Uuid,
		user.Name,
		user.Nickname,
		user.Email,
		hashPass,
		user.UserType,
		user.CreatedAt,
		user.ModifiedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepository) ListUsers() (*[]entities.User, error) {
	query := `
	SELECT u.uuid, u.name, u.nickname, u.email, u.user_type, u.created_at, u.modified_at
	FROM users u
	`

	smtUser, err := ur.db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := smtUser.Close()
		if err != nil {
			return
		}
	}()

	rows, err := smtUser.Query()

	defer func() {
		err := rows.Close()
		if err != nil {
			return
		}
	}()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	users := make([]entities.User, 0)

	for rows.Next() {
		var user entities.User

		err = rows.Scan(
			&user.Uuid,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.UserType,
			&user.CreatedAt,
			&user.ModifiedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}
