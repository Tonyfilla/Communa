package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	login    string `json:"login"`
	password string `json:"password"`
}

func getID(w http.ResponseWriter, ps httprouter.Params) (string, bool) {
	id := ps.ByName("id")
	if id == "" {
		w.WriteHeader(400)
		return "", false
	}
	return id, true
}

func getRecords(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var str string
	if len(r.URL.RawQuery) > 0 {
		str = r.URL.Query().Get("login")
		if str == "" {
			w.WriteHeader(400)
			return
		}
	}
	recs, err := read(str) //тут все верно возвращается
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err = json.NewEncoder(w).Encode(recs); err != nil { // проблема гдето тут
		w.WriteHeader(500)
	}
}

func getRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, ok := getID(w, ps)
	if !ok {
		return
	}
	rec, err := readOne(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err = json.NewEncoder(w).Encode(rec); err != nil {
		w.WriteHeader(500)
	}
}

func addRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var rec User                                // curl -i http://localhost:8080/api/v1/users -XPOST -d '{"login":"IVAN","password":"9284724"}'
	err := json.NewDecoder(r.Body).Decode(&rec) // проблема гдето тут
	if err != nil || rec.password == "" {
		w.WriteHeader(400)
		return
	}
	if _, err := insert(rec.login, rec.password); err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
}

func updateRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, ok := getID(w, ps)
	if !ok {
		return
	}
	var rec User
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil || rec.login == "" || rec.password == "" {
		w.WriteHeader(400)
		return
	}
	res, err := update(id, rec.password)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}

func deleteRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, ok := getID(w, ps)
	if !ok {
		return
	}
	if _, err := remove(id); err != nil {
		w.WriteHeader(500)
	}
	w.WriteHeader(204)
}
