package handlers

import (
	"html/template"
	"net/http"
	"time"

	"work/internal/database"
	"work/internal/models"
	"work/internal/utils"
)

func ClearUserCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:		"user_session",
		Value: 		"",
		Path:		"/",
		MaxAge:		-1,
		Expires:	time.Unix(0, 0),
		HttpOnly: 	true,
	}
	http.SetCookie(w, &cookie)
}


func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login/" {
		http.Redirect(w, r, "/NotFound/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, private, max-age=0, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")

	cookie, err := r.Cookie("user_session")
	if err == nil && cookie.Value != "" {
		var currentUser models.User
		if err := database.DB.Where("name = ?", cookie.Value).First(&currentUser).Error; err == nil {	
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	
	ClearUserCookie(w)

	tmpl, err := template.ParseFiles("public/login.html")
	if err != nil {
		http.Error(w, "HTML file nicht gefunden", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func LoginAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	inputName := r.FormValue("name_input_login")
	inputPassword := r.FormValue("password_input_login")

	var foundUser models.User
	result := database.DB.Where("name = ?", inputName).First(&foundUser)

	if result.Error != nil {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	if !utils.CheckPasswordHash(inputPassword, foundUser.Password) {
		http.Error(w, "Falsches Password", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:		"user_session",
		Value:		foundUser.Name,
		Path: 		"/",
		Expires: 	time.Now().Add(24 * time.Hour),
		HttpOnly: 	true,
		SameSite:	http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register/" {
		http.Redirect(w, r, "/NotFound/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, private, max-age=0, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")

	cookie, err := r.Cookie("user_session")
	if err == nil && cookie.Value != "" {
		var currentUser models.User
		err := database.DB.Where("name = ?", cookie.Value).First(&currentUser)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} 	
	} 

	ClearUserCookie(w)

	tmpl, err := template.ParseFiles("public/register.html")
	if err != nil {
		http.Error(w, "HTML seite nicht gefunden", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)


}

func RegisterAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register/", http.StatusSeeOther)
	}

	inputName := r.FormValue("name_input_register")
	inputPassword := r.FormValue("password_input_register")

	hashedPassword, err := utils.HashPassword(inputPassword)
	if err != nil {
		http.Error(w, "Fehler beim Erstellen den Passwords", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		Name:     inputName,
		Password: hashedPassword,
	}

	result := database.DB.Create(&newUser)
	if result.Error != nil {
		http.Error(w, "Name ist bereits besetzt", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "user_session",
		Value:    newUser.Name,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func LogoutAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ClearUserCookie(w)
	http.Redirect(w, r, "/login/", http.StatusSeeOther)
}