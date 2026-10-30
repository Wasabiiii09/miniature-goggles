package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"work/internal/database"
	"work/internal/handlers"
)

func main() {
	// 1. .env laden
	err := godotenv.Load()
	if err != nil {
		log.Println("Kein .env konfiguriert (ignoriere wenn Umgebungsvariablen gesetzt sind).")
	}

	// 2. Datenbank initialisieren & migrieren
	database.InitDB()

	// 3. Routen registrieren
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/NotFound/", handlers.Page404)
	
	http.HandleFunc("/login/", handlers.LoginPage)
	http.HandleFunc("/login_action", handlers.LoginAction)

	http.HandleFunc("/register/", handlers.RegisterPage)
	http.HandleFunc("/register_action", handlers.RegisterAction)

	http.HandleFunc("/profile/", handlers.ProfilePage)
	http.HandleFunc("/log_out", handlers.LogoutAction)

	http.HandleFunc("/comment/", handlers.CommentPage)
	http.HandleFunc("/comment_action", handlers.CommentAction)
	http.HandleFunc("/delete_comment", handlers.DeleteCommentAction)

	// 4. Server starten
	fmt.Println("Server läuft auf: http://localhost:7779")
	log.Fatal(http.ListenAndServe(":7779", nil))
}