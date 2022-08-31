package entity

import (
	"time"
	"unicode/utf8"
)

const (
	maxLenName   int    = 200
	_defaultName string = "Anonymus"
)

func NewUser(name string) (User, error) {
	if name == "" {
		name = _defaultName
	}

	user := User{Name: name}
	err := user.Check()

	return user, err
}

type User struct {
	Id        int64     `json:"id"`
	CreatedDt time.Time `json:"created_dt"`
	Name      string    `json:"name"`
}

func (u *User) Check() error {
	if utf8.RuneCountInString(u.Name) > maxLenName {
		return ErrorTooManyNameSymbols
	}

	return nil
}
