package main

import (
	"go-htmx/database"
	"go-htmx/handlers"
	"go-htmx/utils"
	"log"
	"net/http"
)

func main() {
	dbConn := database.ConnectDB()
	defer dbConn.Close()

	http.HandleFunc("/hello", handlers.HelloWorldHandler)

	http.HandleFunc("/login", handlers.LoginUser)
	http.HandleFunc("/register", handlers.RegisterUser)

	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/dashboard", handlers.Dashboard)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	logger := utils.NewLogger()
	logger.Info("Server starting...")

	log.Fatal(http.ListenAndServe(":8000", nil))
}
