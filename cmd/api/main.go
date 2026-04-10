package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FranzSinaga/blogcms/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	fmt.Println("Database connected!")

	container := NewContainer(db)
	r := setupRouter(container)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting in %s...", port)
	http.ListenAndServe(":"+port, r)
}
