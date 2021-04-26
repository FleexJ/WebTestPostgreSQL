package main

type UserAndTokens struct {
	User   User
	Tokens []Token
}

//Карта, где хранятся все токены и  авторизированные пользователи
//Ключ - это id пользователя, значение - это пользователь и его токены
type MapTokens map[int]*UserAndTokens

func newMapTokens() *MapTokens {
	result := make(MapTokens)
	return &result
}

//Обновление, передается пользователь, если он есть в карте, то обновляются его данные
func (m *MapTokens) updateUser(u *User) {
	id := u.Id

	//Если его нет в карте, то ничего не делается
	if (*m)[id] == nil {
		return
	}
	(*m)[id].User = *u
}

//Добавление, передается пользователь и токен
func (m *MapTokens) add(u User, t Token) {
	id := u.Id

	//Если этого пользователя нет в карте, создается новая запись с ним
	if (*m)[id] == nil {
		(*m)[id] = &UserAndTokens{
			User:   u,
			Tokens: []Token{t},
		}
		//Если же он уже есть, то к списку его токенов добавляется новый
	} else {
		(*m)[id].User = u
		(*m)[id].Tokens = append((*m)[id].Tokens, t)
	}
}

//Удаляются все токены и сам пользователь из карты
func (m *MapTokens) clearById(id int) {
	delete(*m, id)
}

//Удаляет токен пользователя из карты
func (m *MapTokens) deleteByToken(t Token) {
	id := t.IdUser

	//Если записи нет, то ничего не делаем
	if (*m)[id] == nil {
		return
	}

	//Пересобираем токены без учета удаляемого
	var newSlice []Token
	for _, el := range (*m)[id].Tokens {
		if el.Token != t.Token {
			newSlice = append(newSlice, el)
		}
	}
	(*m)[id].Tokens = newSlice

	//Если не осталось токенов, то удаляем запись в карте
	if len(newSlice) == 0 {
		delete(*m, id)
	}
}

//Получаем хозяина токена, если он есть
func (m MapTokens) getUserByToken(t Token) *User {
	id := t.IdUser

	//Если нет записи этого пользователя
	if m[id] == nil {
		return nil
	}

	//Перебираем все токены пользователя на предмет совпадения, чтобы вернуть искомый
	for _, el := range m[id].Tokens {
		if el.Token == t.Token {
			return &m[id].User
		}
	}
	return nil
}
