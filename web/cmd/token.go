package main

import (
	"encoding/base64"
	"net/http"
	"strconv"
)

const (
	idCookieName    = "id"
	tokenCookieName = "token"
	expDay          = 60 * 24
)

type token struct {
	IdUser int
	Token  string
}

//Сохраняет токен
func (app *application) saveToken(w http.ResponseWriter, u user, t token) {
	http.SetCookie(w,
		newCookie(idCookieName, strconv.Itoa(t.IdUser)))
	//base64 token save in cookie
	base64Tkn := base64.StdEncoding.EncodeToString([]byte(t.Token))
	http.SetCookie(w,
		newCookie(tokenCookieName, base64Tkn))
	app.tokens.add(u, t)
}

//Удаляет токен
func (app *application) deleteToken(w http.ResponseWriter, u user, t token) error {
	http.SetCookie(w,
		newCookie(idCookieName, ""))
	http.SetCookie(w,
		newCookie(tokenCookieName, ""))
	app.tokens.deleteByToken(t)
	return nil
}
