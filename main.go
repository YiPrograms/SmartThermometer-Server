package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var conf config

var db *sql.DB

var queryStmt, registerStmt, delStmt *sql.Stmt

func loadConf(file string) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&conf)
}

func loadDatabase(file string) {
	db, _ = sql.Open("sqlite3", file)
	stmt, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS Cards(
								UID TEXT PRIMARY KEY,
								Num TEXT,
								Name TEXT
							)`)
	stmt.Exec()

	queryStmt, _ = db.Prepare("SELECT Num, Name FROM Cards WHERE UID=?")
	registerStmt, _ = db.Prepare("REPLACE INTO Cards(UID, Num, Name) VALUES(?, ?, ?)")
	delStmt, _ = db.Prepare("DELETE FROM Cards WHERE UID=?")
}

func main() {
	loadConf("config.json")
	loadDatabase("cards.db")

	gsheetsInit()

	http.HandleFunc("/query", query)
	http.HandleFunc("/place", place)
	http.HandleFunc("/register", register)

	log.Fatal(http.ListenAndServeTLS(":9000", "server.crt", "server.key", nil))
}
