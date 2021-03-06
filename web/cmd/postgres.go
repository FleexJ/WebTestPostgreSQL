package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	_ "pq-master"
)

const (
	userConst     = "user=postgres"
	passwordConst = "password=admin"
	dbNameConst   = "dbname=FirstMvc"
	sslmodeConst  = "sslmode=disable"
)

//open connection to DB
func (app application) openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres",
		userConst+" "+
			passwordConst+" "+
			dbNameConst+" "+
			sslmodeConst)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (app application) getUserByEmail(email string) (*User, error) {
	row := app.DB.QueryRow("select * from users where email = $1", email)
	usr := &User{}
	err := row.Scan(&usr.Id, &usr.Name, &usr.Surname, &usr.Email, &usr.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return usr, nil
}

func (app application) getAllUsers() ([]User, error) {
	rows, err := app.DB.Query("select * from users")
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		usr := User{}
		err = rows.Scan(&usr.Id, &usr.Name, &usr.Surname, &usr.Email, &usr.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	return users, nil
}

func (app application) insertUser(usr User) error {
	bcryptPassw, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usr.Password = string(bcryptPassw)
	stmt, err := app.DB.Prepare("insert into users (name, surname, email, password) values ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(usr.Name, usr.Surname, usr.Email, usr.Password)
	if err != nil {
		return err
	}
	return nil
}

func (app application) updateUser(usr User) error {
	stmt, err := app.DB.Prepare("update users set name = $1, surname = $2, email = $3 where id = $4")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(usr.Name, usr.Surname, usr.Email, usr.Id)
	if err != nil {
		return err
	}

	return nil
}

func (app application) updateUserPassword(usr User, password string) error {
	bcryptPassw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := app.DB.Prepare("update users set password = $1 where id = $2")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(string(bcryptPassw), usr.Id)
	if err != nil {
		return err
	}
	return nil
}

func (app application) deleteUser(usr User) error {
	stmt, err := app.DB.Prepare("delete from users where id = $1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(usr.Id)
	if err != nil {
		return err
	}
	return nil
}
