package main

import (
	"encoding/json"
	"net/http"
)

func query(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var qReq queryRequest
	err := json.NewDecoder(r.Body).Decode(&qReq)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if qReq.Key != conf.Key {
		http.Error(w, "Incorrect key", 403)
		return
	}

	var qRes queryResponse

	c, ok := cards[qReq.UID]
	if !ok {
		qRes = queryResponse{-1, ""}
	} else {
		qRes = queryResponse{c.Num, c.Name}
	}

	json.NewEncoder(w).Encode(qRes)
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var rReq registerRequest
	err := json.NewDecoder(r.Body).Decode(&rReq)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if rReq.Key != conf.Key {
		http.Error(w, "Incorrect key", 403)
		return
	}

	cards[rReq.UID] = card{rReq.UID, rReq.Num, ""}
	go saveCards("cards.json")
}

func place(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var pReq placeRequest
	err := json.NewDecoder(r.Body).Decode(&pReq)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if pReq.Key != conf.Key {
		http.Error(w, "Incorrect key", 403)
		return
	}

	go writeTemp(pReq.Num, pReq.Temp)
}