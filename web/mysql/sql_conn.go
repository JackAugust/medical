package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Initdb() *sql.DB {
	db, _ = sql.Open("mysql", "root:root@tcp(localhost:3306)/itbtsql?charset=utf8&allowNativePasswords=true")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect mysql , err :" + err.Error())
		os.Exit(1)
	} else {
		//fmt.Println("连接成功")
	}
	return db
}
