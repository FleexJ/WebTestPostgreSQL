package main

import (
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)


func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Expires: time.Now().Add(expDay * time.Hour),
	}
}

//Возвращает токен, считанный из куки
func (app *application) getTokenCookies(r *http.Request) *token {
	cookieId, err := r.Cookie(idCookieName)
	if err != nil {
		return nil
	}
	cookieToken, err := r.Cookie(tokenCookieName)
	if err != nil {
		return nil
	}
	if cookieId.Value == "" || cookieToken.Value == "" {
		return nil
	}
	id, err := strconv.Atoi(cookieId.Value)
	if err != nil {
		return nil
	}
	return &token{
		IdUser: id,
		Token:  cookieToken.Value,
	}
}

//Функция авторизации пользователя
//Ищет совпадения в базе пользователей
//Выдает новый токен доступа
//при успехе нет ошибки
func (app *application) auth(w http.ResponseWriter, email, password string) error {
	u := app.getUserByEmail(email)
	if u == nil {
		return errors.New("user not found")
	}
	err := u.comparePassword(password)
	if err != nil {
		return err
	}
	genToken, err := app.generateToken(strconv.Itoa(u.Id))
	if err != nil {
		return err
	}
	tkn := token{
		IdUser: u.Id,
		Token:  genToken,
	}
	app.saveToken(w, *u, tkn)
	return nil
}

//Проверка токена доступа, возвращает токен с данными и текущего пользователя при успехе
func (app *application) checkAuth(r *http.Request) (*token, *user) {
	tkn := app.getTokenCookies(r)
	if tkn == nil {
		return nil, nil
	}
	//Декодируем токен из куки
	tDecode, err := base64.StdEncoding.DecodeString(tkn.Token)
	if err != nil {
		return nil, nil
	}
	tkn.Token = string(tDecode)
	u := app.tokens.getUserByToken(*tkn)
	if u == nil {
		return nil, nil
	}
	return tkn, u
}

//Генерирует новый токен на основе слова
func (app *application) generateToken(word string) (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := 20
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	bcryptB, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return word + strconv.FormatInt(time.Now().Unix(), 10) + string(bcryptB), nil
}
