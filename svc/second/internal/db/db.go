package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

var db *sql.DB

func Init(user, pass, host, port, name string) func() {
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	sqltrace.Register("mysql", &mysql.MySQLDriver{})
	var err error
	db, err = sqltrace.Open("mysql", dbString)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open database connection: %w", err))
	}
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Errorf("failed to ping database: %w", err))
	}

	db.SetConnMaxLifetime(time.Minute * 15)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)

	initDB()

	return func() { _ = db.Close() }
}

//go:embed schema.sql
var schema string

func initDB() {
	stmts := strings.Split(schema, ";")
	for _, stmt := range stmts[:len(stmts)-1] {
		stmt = strings.Trim(stmt, "\n ") + ";"
		if _, err := db.Exec(stmt); err != nil {
			log.Fatal(err)
		}
	}
}
