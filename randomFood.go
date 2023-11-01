package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

func randomFood(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://www.themealdb.com/api/json/v1/1/random.php")

	if err != nil {
		fmt.Printf("Error getting random food: %s\n", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	var data FoodApiTypes
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("Error unmarshalling response body: %s\n", err)
		return
	}
	templ, err := template.ParseFiles("./templates/foodTemplate.html")
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		return
	}

	var tpl bytes.Buffer
	err = templ.Execute(&tpl, data.Meals[0])
	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(tpl.Bytes())
}