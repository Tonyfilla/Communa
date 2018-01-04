package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/user"

	"github.com/FogCreek/mini"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func params() string {
	u, err := user.Current()
	fatal(err)
	cfg, err := mini.LoadConfiguration(u.HomeDir + "/go/src/webserver_communa/.communarc")
	fatal(err)

	info := fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		cfg.String("host", "127.0.0.1"),
		cfg.String("port", "5432"),
		cfg.String("dbname", u.Username),
		cfg.String("sslmode", "disable"),
		cfg.String("user", u.Username),
		cfg.String("pass", ""),
	)
	return info
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "/Users/tony/go/src/webserver_communa/communa.db")
	fatal(err)
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
		`users("login" TEXT, "pasword" TEXT)`) //какие типы лучше?
	fatal(err)

	router := httprouter.New()
	router.GET("/api/v1/users", getRecords)
	router.GET("/api/v1/users/:id", getRecord)
	router.POST("/api/v1/users", addRecord)
	router.PUT("/api/v1/users/:id", updateRecord)
	router.DELETE("/api/v1/users/:id", deleteRecord)
	http.ListenAndServe(":8080", router)
}
