package easyWeb

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DBname = "database.db"

func CreateDB(sqlpath string) {
	fileExists := func(filename string) bool {
		info, err := os.Stat(DBname)
		if os.IsNotExist(err) {
			return false
		}
		return !info.IsDir()
	}

	if !fileExists(DBname) {
		query, errorRead := ioutil.ReadFile(sqlpath)
		err(errorRead)

		db, errorOpen := sql.Open("sqlite3", DBname)
		err(errorOpen)
		defer db.Close()
		_, errorExec := db.Exec(string(query))
		err(errorExec)
	}
}

func Insert(query string, args ...interface{}) error {
	db, errorOpen := sql.Open("sqlite3", DBname)
	err(errorOpen)
	defer db.Close()
	tx, errorBegin := db.Begin()
	err(errorBegin)
	_, errorExec := tx.Exec(query, args...)
	if errorExec != nil {
		err(tx.Rollback())
		return errorExec
	}
	err(tx.Commit())
	return nil
}
