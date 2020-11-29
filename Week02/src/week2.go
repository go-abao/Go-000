package main

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// DB sql.DB
var DB *sql.DB

const userTable = "`user`"

// User User
type User struct {
	ID       int    `json:"id"`       // ID
	Username string `json:"username"` // 用户名
}

func main() {
	var id = 2
	userInfo, err := User{ID: id}.GetUsernameByID()
	if err != nil {
		log.Errorf("%+v\n", err)
	}
	log.Println(userInfo)
}

// GetUsernameByID get username by id
func (u User) GetUsernameByID() (User, error) {
	var user User
	err := sq.
		Select("id, username").
		From(userTable).
		Where(sq.Eq{"id": u.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&user.ID, &user.Username)
	if err == sql.ErrNoRows {
		return user, errors.New("Record does not exist")
	} else if err != nil {
		return user, err
	}
	return user, nil
}
