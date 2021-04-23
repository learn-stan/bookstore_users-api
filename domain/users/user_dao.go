package users

import (
	"fmt"

	"github.com/learn-stan/bookstore_users-api/datasources/mysql/users_db"
	"github.com/learn-stan/bookstore_users-api/logger"
	"github.com/learn-stan/bookstore_users-api/utils/errors"
)

const (
	errorNoRows           = "no rows in result set"
	indexUniqueEmail      = "email_UNIQUE"
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
)

func (user *User) Get() *errors.RestErr {
	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getError := result.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.DateCreated,
		&user.Status,
	); getError != nil {
		logger.Error("error when trying to prepare get user by id", getError)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare set user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	insertResult, saveError := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveError != nil {
		logger.Error("error when trying to save user", saveError)
		return errors.NewInternalServerError("database error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerError("database error")
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update query", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user query", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare search user query", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to search user", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	result := make([]User, 0)

	for rows.Next() {
		var user User

		if getError := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.DateCreated,
			&user.Status,
		); getError != nil {
			logger.Error("error when trying to scan user row into user struct", getError)
			return nil, errors.NewInternalServerError("database error")
		}

		result = append(result, user)
	}

	if len(result) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return result, nil
}
