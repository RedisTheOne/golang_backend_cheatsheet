package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db, dbErr = sql.Open("mysql", "root:toor@tcp(127.0.0.1:3306)/login")
