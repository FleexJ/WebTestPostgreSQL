package main

import (
	"golang.org/x/crypto/bcrypt"
)

const regexEmail = `^\w+@\w+[.]\w+$`

type User struct {
	Id       int
	Name     string
	Surname  string
	Email    string
	Password string
}

func (u *User) comparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == nil {
		return nil
	}
	return err
}
