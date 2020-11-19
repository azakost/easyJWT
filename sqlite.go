package easyWeb

import (
	"database/sql"
	"io/ioutil"
	"os"
	"reflect"

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

func Exe(query string, args ...interface{}) {
	db, errorOpen := sql.Open("sqlite3", DBname)
	err(errorOpen)
	defer db.Close()
	tx, errorBegin := db.Begin()
	err(errorBegin)
	_, errorExec := tx.Exec(query, args...)
	if errorExec != nil {
		err(tx.Rollback())
		panic(errorExec)
	}
	err(tx.Commit())
}

func Select(model interface{}, query string, args ...interface{}) {
	db, databaseError := sql.Open("sqlite3", DBname)
	err(databaseError)
	defer db.Close()
	statement, stmError := db.Prepare(query)
	err(stmError)
	rows, rowsError := statement.Query(args...)
	err(rowsError)
	defer rows.Close()
	container := reflect.Indirect(reflect.ValueOf(model))
	v := container.Type().Elem()
	len := v.NumField()
	tmp := make([]interface{}, len)
	dest := make([]interface{}, len)
	for i := range tmp {
		dest[i] = &tmp[i]
	}
	for rows.Next() {
		scanError := rows.Scan(dest...)
		err(scanError)
		row := reflect.Indirect(reflect.New(v))
		for i, t := range tmp {
			f := row.Field(i)
			f.Set(reflect.ValueOf(t))
		}
		container.Set(reflect.Append(container, row))
	}
}

func InDB(query string, data interface{}) bool {
	db, databaseError := sql.Open("sqlite3", DBname)
	err(databaseError)
	defer db.Close()
	queryError := db.QueryRow(query, data).Scan(&data)
	if queryError != nil {
		if queryError != sql.ErrNoRows {
			panic(queryError)
		}
		return false
	}
	return true
}
