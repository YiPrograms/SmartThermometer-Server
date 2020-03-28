package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func query(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	log.Println(ip, ": Query request from", ip)

	r.ParseForm()
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		log.Println(ip, ": 400 No request body")
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(ip, ": 500 Parse body error:", err)
		return
	}

	body := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Println(ip, ":", body)

	var qReq queryRequest
	err = json.NewDecoder(r.Body).Decode(&qReq)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Println(ip, ": 400 JSON decode error:", err)
		return
	}

	if qReq.Key != conf.Key {
		http.Error(w, "Incorrect key", 403)
		log.Println(ip, ": 403 Incorrect key:", qReq.Key)
		return
	}

	var Num int
	var Name string

	err = queryStmt.QueryRow(qReq.UID).Scan(&Num, &Name)

	var qRes queryResponse

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No entry found", 422)
			log.Println(ip, ": 422 No entry found")
		} else {
			http.Error(w, "SQL Error", 500)
			log.Println(ip, ": 500 SQL Error:", err)
		}
		return
	}
	qRes = queryResponse{Num, Name}
	json.NewEncoder(w).Encode(qRes)
	log.Println(ip, ": 200 OK")
}

func register(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	log.Println(ip, ": Register request from", ip)

	r.ParseForm()
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		log.Println(ip, ": 400 No request body")
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(ip, ": 500 Parse body error:", err)
		return
	}

	body := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Println(ip, ":", body)

	var rReq registerRequest
	err = json.NewDecoder(r.Body).Decode(&rReq)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Println(ip, ": 400 JSON decode error:", err)
		return
	}

	if rReq.Key != conf.Key {
		http.Error(w, "Incorrect key", 403)
		log.Println(ip, ": 403 Incorrect key:", rReq.Key)
		return
	}

	_, err = registerStmt.Exec(rReq.UID, rReq.Num, rReq.Name)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "SQL Error", 500)
		log.Println(ip, ": 500 SQL Error:", err)
		return
	}

	log.Println(ip, ": 200 OK")
}

func deregister(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	log.Println(ip, ": Deregister request from", ip)

	r.ParseForm()
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		log.Println(ip, ": 400 No request body")
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(ip, ": 500 Parse body error:", err)
		return
	}

	body := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Println(ip, ":", body)

	var dReq deregisterRequest
	err = json.NewDecoder(r.Body).Decode(&dReq)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Println(ip, ": 400 JSON decode error:", err)
		return
	}

	if dReq.Key != conf.Key {
		http.Error(w, "Incorrect key", 403)
		log.Println(ip, ": 403 Incorrect key:", dReq.Key)
		return
	}

	res, err := delStmt.Exec(dReq.UID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "SQL Error", 500)
		log.Println(ip, ": 500 SQL Error:", err)
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Not registered", 422)
		log.Println(ip, ": 422 Not registered")
		return
	}

	log.Println(ip, ": 200 OK")
}

func place(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	log.Println(ip, ": Place request from", ip)

	r.ParseForm()
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		log.Println(ip, ": 400 No request body")
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(ip, ": 500 Parse body error:", err)
		return
	}

	body := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Println(ip, ":", body)

	var pReq placeRequest
	err = json.NewDecoder(r.Body).Decode(&pReq)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Println(ip, ": 400 JSON decode error:", err)
		return
	}

	if pReq.Key != conf.Key {
		http.Error(w, "Incorrect key", 403)
		log.Println(ip, ": 403 Incorrect key:", pReq.Key)
		return
	}

	go writeTemp(pReq.Num, pReq.Temp)

	w.WriteHeader(202)
	log.Println(ip, ": 202 Accepted")
}
