package main

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const regexEmail = `^\w+@\w+[.]\w+$`

type User struct {
	Id       int
	Name     string
	Surname  string
	Email    string
	Password string
}

func (usr *User) comparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err == nil {
		return nil
	}
	return err
}

func (app application) validUser(usr User, repPassword string) (bool, error) {
	matched, _ := regexp.MatchString(regexEmail, usr.Email)
	if !matched ||
		usr.Name == "" ||
		usr.Surname == "" ||
		usr.Password == "" ||
		usr.Password != repPassword {
		return false, nil
	}
	uEmail, err := app.getUserByEmail(usr.Email)
	if err != nil {
		return false, err
	}
	if uEmail != nil && uEmail.Id != usr.Id {
		return false, nil
	}
	return true, nil
}
