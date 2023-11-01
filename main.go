package main

import (
	"database/sql"
	"errors"
	"fmt"

	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {fmt.Printf("slaying & serving index.html \n"); http.ServeFile(w, r, "./index.html")})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {fmt.Printf("slaying & serving about.html\n"); http.ServeFile(w, r, "./about.html")})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {fmt.Printf("slaying & serving contact.html\n"); http.ServeFile(w, r, "./contact.html")})
	http.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {fmt.Printf("slaying & serving random.html\n"); http.ServeFile(w, r, "./random.html")})
	http.HandleFunc("/test", testRoute)
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/selectAll", selectAll)
	http.HandleFunc("/randomFood", randomFood)

	fmt.Print("\033[H\033[2J")
	text := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, "\n\n\n ** Server Started at http://localhost:3333 ** \n\n\n\n\n\n\n\n\n\n\n\n")
	fmt.Print(text)
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

func testRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /test request\n")
	io.WriteString(w, "This is server side data\n")
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



for rows.Next() {
	user := new(Users)
	err = rows.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		fmt.Println(err)
	}
	tableRow := "<script> alert('" + fmt.Sprint(user.ID) + ". " + user.Name + user.Email + "')</script>"

	io.WriteString(w, tableRow)
}
db.Close()
}



