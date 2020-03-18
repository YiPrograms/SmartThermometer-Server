package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var conf config
var cards = make(map[string]card)

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

func main() {
	loadConf("config.json")
	loadCards("cards.json")
	defer saveCards("cards.json")

	gsheetsInit()

	http.HandleFunc("/query", query)
	http.HandleFunc("/place", place)
	http.HandleFunc("/register", register)

	log.Fatal(http.ListenAndServeTLS(":9000", "server.crt", "server.key", nil))
}
