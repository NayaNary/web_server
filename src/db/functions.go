package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Db struct {
	obj       *sql.DB
	lastError error
}

func NewConnect(nameDriver string, connParam string) *Db {
	db, err := sql.Open(nameDriver, connParam)
	return &Db{obj: db, lastError: err}
}

func (db *Db) Close() {
	db.lastError = db.obj.Close()
}

func (db *Db) LastError() error {
	return db.lastError
}

func (db *Db) Insert(table string, rows []string, values []interface{}) (result sql.Result) {
	db.lastError = nil
	if len(table) == 0 {
		db.lastError = errors.New("не передано название таблицы")
		return
	}
	if len(rows) == 0 {
		db.lastError = errors.New("не передано ни одного поля для записи")
		return
	}
	if len(values) == 0 {
		db.lastError = errors.New("не передано ни одного значения для записи")
		return
	}

	sqlQuery := "insert into " + table

	sqlQuery += addRows(rows, "insert")

	sqlQuery += "values "

	sqlQuery += addValues(values, "insert")
	result, db.lastError = db.obj.Exec(sqlQuery)
	return
}

func addRows(data []string, queryName string) (table string) {
	if len(data) != 0 {
		if strings.ToLower(queryName) == "insert" {
			table = "("
		}

		for _, row := range data {
			table += row + ","
		}

		table = strings.TrimRight(table, ",")

		if strings.ToLower(queryName) == "insert" {
			table += ")"
		}

	}

	return
}

func addValues(data []interface{}, queryName string) (values string) {
	if len(data) != 0 {
		if strings.ToLower(queryName) == "insert" {
			values = "("
		}

		for _, value := range data {
			typeValue := fmt.Sprintf("%T", value)
			fmt.Println("type:",typeValue)
		
			if typeValue == "string" || typeValue == "[]string" {
				values += fmt.Sprintf("'%v',", value)
			} else {
				values += fmt.Sprintf("%v,", value)
			}

		}

		values = strings.TrimRight(values, ",")

		if strings.ToLower(queryName) == "insert" {
			values += ")"
		}

	}

	return
}

func (db *Db) Select(table string, rows []string, where string) (result *sql.Rows) {
	db.lastError = nil
	if len(table) == 0 {
		return
	}
	if len(rows) == 0 {
		return
	}

	sqlQuery := "Select " + addRows(rows, "select") + " from " + table + " " + where
	fmt.Println(sqlQuery)
	result, db.lastError = db.obj.Query(sqlQuery)
	fmt.Println("err", db.lastError)
	return
}
