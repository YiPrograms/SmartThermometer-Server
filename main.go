package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var conf config
var cards = make(map[string]card)

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

func loadCards(file string) {
	cardsFile, err := os.Open(file)
	defer cardsFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	var cardsJSON []card
	json.NewDecoder(cardsFile).Decode(&cardsJSON)

	for _, c := range cardsJSON {
		cards[c.UID] = c
	}
}

func saveCards(file string) {
	cardsFile, err := os.Open(file)
	defer cardsFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	var cardsArray []card
	for _, c := range cards {
		cardsArray = append(cardsArray, c)
	}

	JSONString, _ := json.MarshalIndent(cardsArray, "", " ")
	ioutil.WriteFile(file, JSONString, os.ModePerm)
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
	loadCards("cards.json")
	defer saveCards("cards.json")

	gsheetsInit()

	http.HandleFunc("/query", query)
	http.HandleFunc("/place", place)
	http.HandleFunc("/register", register)

	log.Fatal(http.ListenAndServeTLS(":9000", "server.crt", "server.key", nil))
}
