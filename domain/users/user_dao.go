package users

import (
	"fmt"

	"github.com/learn-stan/bookstore_users-api/datasources/mysql/users_db"
	"github.com/learn-stan/bookstore_users-api/utils/date_utils"
	"github.com/learn-stan/bookstore_users-api/utils/errors"
	"github.com/learn-stan/bookstore_users-api/utils/mysql_utils"
)

const (
	errorNoRows      = "no rows in result set"
	indexUniqueEmail = "email_UNIQUE"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?, ?, ?, ?);"
	queryGetUser     = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
	queryUpdateUser  = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
)

func (user *User) Get() *errors.RestErr {
	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewBadRequestError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getError := result.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.DateCreated,
	); getError != nil {
		return mysql_utils.ParseError(getError)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	insertResult, saveError := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if saveError != nil {
		return mysql_utils.ParseError(saveError)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error, when trying to save user: %s", err.Error()))
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}
