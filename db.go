package main

import (
	"database/sql"
)

func insert(login, password string) (sql.Result, error) {
	return db.Exec("INSERT INTO users VALUES ($1, $2)",
		login, password)
}

func remove(login string) (sql.Result, error) {
	return db.Exec("DELETE FROM users WHERE login=$1", login)
}

func update(login, password string) (sql.Result, error) {
	return db.Exec("UPDATE users SET password = $1 WHERE login=$2",
		password, login)
}

func readOne(login string) (User, error) {
	var rec User
	row := db.QueryRow("SELECT * FROM users WHERE login=$1 ORDER BY login", login)
	return rec, row.Scan(&rec.Login, &rec.Password) //how it work
}

func read(str string) ([]User, error) {
	var rows *sql.Rows
	var err error
	if str != "" {
		rows, err = db.Query("SELECT * FROM users WHERE login LIKE $1 ORDER BY id",
			"%"+str+"%")
	} else {
		rows, err = db.Query("SELECT * FROM users ORDER BY login")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rs = make([]User, 0)
	var rec User
	for rows.Next() {
		if err = rows.Scan(&rec.Login, &rec.Password); err != nil {
			return nil, err
		}
		rs = append(rs, rec)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rs, nil
}
