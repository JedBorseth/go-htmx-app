package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/mrz1836/go-sanitize"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {fmt.Printf("req / \n"); http.ServeFile(w, r, "./index.html")})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {fmt.Printf("got /about request\n"); http.ServeFile(w, r, "./about.html")})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {fmt.Printf("got /contact request\n"); http.ServeFile(w, r, "./contact.html")})
	http.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {fmt.Printf("got /random request\n"); http.ServeFile(w, r, "./random.html")})
	http.HandleFunc("/test", testRoute)
	http.HandleFunc("/hello", helloRoute)
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/selectAll", selectAll)
	http.HandleFunc("/randomFood", randomFood)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		log.Fatal(err)
	}
}
func getEnv(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }
  func InitDb() *sql.DB {
	dbStr := getEnv("PLANETSCALE_URL")
	if (len(dbStr) <= 0) {
		log.Fatal("No DB String Found")
	}
	db, err := sql.Open("mysql", dbStr)
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}
	
		if err := db.Ping(); err != nil {
			log.Fatalf("failed to ping: %v", err)
		}
	
		log.Println("Successfully connected to PlanetScale!")
		return db
}

func helloRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func testRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /test request\n")
	io.WriteString(w, "This is server side data\n")
}

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

type Users struct {
	ID            int64  `field:"id"`                      
    Name      string `field:"name"`           
    Email      string `field:"email"`                     
}
func selectAll(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /selectAll request\n")
	w.Header().Set("Content-Type", "application/json")

	db := InitDb()
	queryStr := "SELECT * FROM users;"
	rows, err := db.Query(queryStr)
	if(err != nil) {
		log.Fatal("Failed to get users from DB")
	}


	// list := []string{}
for rows.Next() {
	user := new(Users)
	err = rows.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		fmt.Println(err)
	}
	tableRow := "<tr class='w-full flex gap-2'><td class='lg:w-16'>" + fmt.Sprint(user.ID) + "</td><td>" + user.Name + "</td><td>" + user.Email + "</td></tr>"
	// list = append(list, tableRow)
	io.WriteString(w, tableRow)
}
// jsonResp, err := json.Marshal(list)
// 		if(err != nil) {
// 			log.Fatal(err)
// 		}
// w.Write(jsonResp)

}


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
