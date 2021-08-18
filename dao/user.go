package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/muchlist/sagasql/db"
	"github.com/muchlist/sagasql/dto"
	"github.com/muchlist/sagasql/utils"
)

func NewUserDao() UserDaoAssumer {
	return &userDao{}
}

type UserDaoAssumer interface {
	Insert(user dto.User) (*int64, rest_err.APIError)
	Get(id int64) (*dto.User, rest_err.APIError)
	Find() ([]dto.User, rest_err.APIError)
	Edit(userInput dto.User) (*dto.User, rest_err.APIError)
	Delete(userID int64) rest_err.APIError
}

type userDao struct {
}

func (u *userDao) Insert(user dto.User) (*int64, rest_err.APIError) {
	sqlStatement := `
	INSERT INTO users (username, email, name, password, roles, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id;
	`
	var userID int64
	err := db.DB.QueryRow(context.Background(), sqlStatement, user.Username, user.Email, user.Name, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&userID)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal menambahkan user", err)
	}
	return &userID, nil
}

func (u *userDao) Edit(input dto.User) (*dto.User, rest_err.APIError) {
	sqlStatement := `
	UPDATE users 
	SET username = $2, email = $3, name = $4, role = $5, updated_at = $6
	WHERE user_id = $1 
	RETURNING user_id, username, email, name, role, created_at, updated_at;
	`

	var user dto.User
	err := db.DB.QueryRow(
		context.Background(),
		sqlStatement, input.UserID, input.Username, input.Email, input.Name, input.Role, input.UpdatedAt,
	).Scan(&user.UserID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("User dengan ID %d tidak ditemukan", input.UserID))
		} else {
			return nil, rest_err.NewInternalServerError("gagal mengedit user", err)
		}
	}
	return &user, nil
}

func (u *userDao) ChangePassword(input dto.User) (*dto.User, rest_err.APIError) {
	sqlStatement := `
	UPDATE users 
	SET hash_pw = $2, updated_at = $3 
	WHERE user_id = $1 
	RETURNING user_id, username, email, name, role, created_at, updated_at;
	`

	var user dto.User
	err := db.DB.QueryRow(context.Background(), sqlStatement, input.UserID, input.Password, input.UpdatedAt).Scan(&user.UserID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("User dengan ID %d tidak ditemukan", input.UserID))
		} else {
			return nil, rest_err.NewInternalServerError("gagal mengganti password user", err)
		}
	}
	return &user, nil
}

func (u *userDao) Delete(userID int64) rest_err.APIError {
	sqlStatement := `
	DELETE FROM users 
	WHERE user_id = $1;
	`
	res, err := db.DB.Exec(context.Background(), sqlStatement, userID)
	if err != nil {
		return rest_err.NewInternalServerError("gagal saat penghapusan user", err)
	}
	if res.RowsAffected() != 1 {
		return rest_err.NewBadRequestError(fmt.Sprintf("User dengan ID %d tidak ditemukan", userID))
	}

	return nil
}

func (u *userDao) Get(id int64) (*dto.User, rest_err.APIError) {

	sqlStatement := `
	SELECT user_id, username, email, name, password, role, created_at, updated_at 
	FROM users 
	WHERE user_id = $1;
	`
	row := db.DB.QueryRow(context.Background(), sqlStatement, id)

	var user dto.User
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Name, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("User dengan ID %d tidak ditemukan", id))
		} else {
			return nil, rest_err.NewInternalServerError("gagal mendapatkan detil user", err)
		}
	}
	return &user, nil
}

func (u *userDao) Find() ([]dto.User, rest_err.APIError) {
	rows, err := db.DB.Query(context.Background(),
		"SELECT user_id, username, email, name, role, created_at, updated_at  FROM users;")
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar user", err)
	}

	var users []dto.User
	for rows.Next() {
		user := dto.User{}
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, rest_err.NewInternalServerError("gagal scan list user", err)
		}
		users = append(users, user)
	}
	defer rows.Close()
	return users, nil
}
