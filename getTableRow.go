package main

import (
	"fmt"
	"io"
	"net/http"
)

func getTableRow(w http.ResponseWriter, r *http.Request)  {
	id := 1
	db := InitDb()
	queryString := "SELECT * FROM users WHERE id = '" + fmt.Sprint(id) + "';"
	rows, err := db.Query(queryString)
	if(err != nil){
		panic(err.Error())
	}
	db.Close()
	
	for rows.Next() {
		user := new(Users)
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			fmt.Println(err)
		}
	io.WriteString(w, "test")
	}
}