package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
)

func getTableRow(w http.ResponseWriter, r *http.Request)  {
	current := 0
	db := InitDb()
	queryString := "SELECT * FROM users WHERE id > " + fmt.Sprint(current) + ";"
	rows, err := db.Query(queryString)
	if(err != nil){
		fmt.Print("Error querying DB\n", err, "\n")
		return
	}
	io.WriteString(w, "<table class='divide-y divide-gray-300 border rounded'><thead class='bg-gray-50'><tr><th class='px-6 py-2 text-xs text-gray-500'>ID</th><th class='px-6 py-2 text-xs text-gray-500'>Name</th><th class='px-6 py-2 text-xs text-gray-500'>Email</th><th class='px-6 py-2 text-xs text-gray-500'>Updated at</th><th class='px-6 py-2 text-xs text-gray-500'>Edit</th><th class='px-6 py-2 text-xs text-gray-500'>Delete</th></tr></thead><tbody class='bg-white divide-y divide-gray-300'>")
	for rows.Next() {
		user := new(Users)
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Print("Error scanning rows\n", err, "\n")
			return
		}
		templ, err := template.ParseFiles("templates/tableRowTemplate.html")
		if(err != nil) {
		fmt.Print("Error parsing row template\n ", err)
		return
		}
		templ.Execute(w, user)
	}
	io.WriteString(w, "</tbody></table>")
	db.Close()
	
}
