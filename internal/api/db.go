package api

import (
	"database/sql"

	"log"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func initDB() *DB {
	db, err := sql.Open("sqlite", config.getDbPath())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
	create table IF NOT EXISTS clients (client_id text not null primary key, client_secret text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{db}
}

func (db *DB) ClearWholeDB() {
	sqlStmt := "delete from clients;"
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}
