package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mrz1836/go-sanitize"
)

func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /addUser request\n")

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)

	// parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		resp["error"] = "Failed to parse form data"
		return
	}

	// get form data
	name := sanitize.Alpha(r.Form.Get("name"), true)
	email := sanitize.Email(r.Form.Get("email"), false)

	// Sanitize and check length of fields
	if(len(name) <= 0 || len(email) <= 0) {
		http.Error(w, "Missing Required Fields", http.StatusBadRequest)
		resp["error"] = "Missing Required Fields"
		return
	}
	// Init db
		db := InitDb()
		execStr := "INSERT INTO users (name, email) values ('" + name +  "', '" + email + "')"
		res, err := db.Exec(execStr)
		if(err != nil) {
			log.Fatalf("err: %s\n", err)
		}
		fmt.Print(res, "\n")
		db.Close()
	// Send Json Response
		fmt.Print("User added successfully\n")
		w.WriteHeader(http.StatusOK)
		resp["success"] = "Added User"
		jsonResp, err := json.Marshal(resp)
		if(err != nil) {
			log.Fatal(err)
		}
		w.Write(jsonResp)

}
