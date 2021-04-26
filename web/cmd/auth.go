package main

import (
	"encoding/base64"
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
func (app *application) getTokenCookies(r *http.Request) *Token {
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

	return &Token{
		IdUser: id,
		Token:  cookieToken.Value,
	}
}

//Проверка токена доступа, возвращает токен с данными и текущего пользователя при успехе
func (app *application) checkAuth(r *http.Request) (*Token, *User) {
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
