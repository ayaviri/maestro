package db

import (
	"database/sql"
	"errors"
	"fmt"
	"maestro/internal"

	"github.com/google/uuid"
)

var err error

type User struct {
	Id       int64
	Username string
	Email    string
}

func GetUserFromToken(db *sql.DB, token string) (User, error) {
	query := fmt.Sprintf(
		`select id, username, email from user where token = "%s"`, token,
	)
	var row *sql.Row = db.QueryRow(query)
	var user User
	err = row.Scan(&user.Id, &user.Username, &user.Email)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func IsValidToken(db *sql.DB, token string) (bool, error) {
	var count int64
	query := fmt.Sprintf(
		"select count(*) from user where token=\"%s\";", token,
	)
	count, err = QueryCount(db, query)

	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// Returns true if there are any users in the `user` table with
// the given username. False otherwise
func IsUsernameAvailable(db *sql.DB, username string) (bool, error) {
	var count int64
	query := fmt.Sprintf(
		"select count(*) from user where username=\"%s\";", username,
	)
	count, err = QueryCount(db, query)

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func CreateUser(
	db *sql.DB,
	username string,
	password string,
	email string,
) (int64, error) {
	var passwordHash string
	passwordHash, err = internal.HashString(password)

	if err != nil {
		return 0, err
	}

	statement := fmt.Sprintf(
		`insert into user (username, password_hash, email)
values("%s", "%s", "%s");`,
		username, passwordHash, email,
	)
	var result sql.Result
	result, err = db.Exec(statement)

	if err != nil {
		return 0, err
	}

	var userId int64
	userId, err = result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return userId, nil
}

// 1) Checks for the existence of a user with the given credentials
// 2) Creates, writes to the DB, and returns a UUID bearer token
// if there exists exactly one user with the given credentials
func AuthenticateAndGenerateToken(
	db *sql.DB,
	username string,
	password string,
) (string, error) {
	var passwordHash string
	passwordHash, err = internal.HashString(password)

	if err != nil {
		return "", err
	}

	var count int64
	query := fmt.Sprintf(
		`select count(*) from user where username="%s" and password_hash="%s";`,
		username, passwordHash,
	)
	count, err = QueryCount(db, query)

	if err != nil {
		return "", err
	}

	if count != 1 {
		return "", errors.New("Invalid credentials")
	}

	bearerToken := uuid.NewString()
	err = UpdateBearerToken(db, username, bearerToken)

	if err != nil {
		return "", err
	}

	return bearerToken, nil
}

func UpdateBearerToken(db *sql.DB, username, bearerToken string) error {
	statement := fmt.Sprintf(
		"update user set token = \"%s\" where username = \"%s\"",
		bearerToken, username,
	)
	var result sql.Result
	result, err = db.Exec(statement)

	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()

	if err != nil {
		return nil
	}

	if rowsAffected != 1 {
		// TODO: How do we roll this back...
		return fmt.Errorf(
			"%d rows affected during token update", rowsAffected,
		)
	}

	return nil
}
