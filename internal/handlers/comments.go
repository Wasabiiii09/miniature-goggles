package handlers

import (
	"html/template"
	"net/http"

	"work/internal/database"
	"work/internal/models"
)

// GET /NotFound/
func Page404(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("public/404.html")
	if err != nil {
		http.Error(w, "HTML not found", http.StatusSeeOther)
		return
	}
	tmpl.Execute(w, nil)
}

// GET /
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/NotFound/", http.StatusSeeOther)
		return
	}

	var currentUser models.User
	var isLoggedIn bool

	cookie, err := r.Cookie("user_session")
	if err == nil && cookie.Value != "" {
		result := database.DB.Where("name = ?", cookie.Value).First(&currentUser)
		if result.Error == nil {
			isLoggedIn = true
		}
	}

	var allComments []models.Comment
	database.DB.Preload("User").Order("created_at desc").Find(&allComments)

	data := struct {
		User       models.User
		Comments   []models.Comment
		IsLoggedIn bool
	}{
		User:       currentUser,
		Comments:   allComments,
		IsLoggedIn: isLoggedIn,
	}

	homeTmpl, err := template.ParseFiles("public/home.html")
	if err != nil {
		http.Error(w, "HTML nicht gefunden", http.StatusInternalServerError)
		return
	}
	homeTmpl.Execute(w, data)
}

// GET /profile/
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile/" {
		http.Redirect(w, r, "/NotFound/", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("user_session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	var currentUser models.User
	result := database.DB.Where("name = ?", cookie.Value).First(&currentUser)
	if result.Error != nil {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	var userComments []models.Comment
	database.DB.Where("user_id = ?", currentUser.ID).Order("created_at desc").Find(&userComments)

	data := struct {
		User     models.User
		Comments []models.Comment
	}{
		User:     currentUser,
		Comments: userComments,
	}

	profileTmpl, err := template.ParseFiles("public/profile.html")
	if err != nil {
		http.Error(w, "HTML nicht gefunden", http.StatusInternalServerError)
		return
	}
	profileTmpl.Execute(w, data)
}

// POST /delete_comment
func DeleteCommentAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile/", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("user_session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	var currentUser models.User
	if err := database.DB.Where("name = ?", cookie.Value).First(&currentUser).Error; err != nil {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	commentID := r.FormValue("comment_id")

	// Nur eigene Kommentare löschen
	result := database.DB.Where("id = ? AND user_id = ?", commentID, currentUser.ID).Delete(&models.Comment{})
	if result.Error != nil {
		http.Error(w, "Fehler beim Löschen des Kommentars", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile/", http.StatusSeeOther)
}

// GET /comment/
func CommentPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/" {
		http.Redirect(w, r, "/NotFound/", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("user_session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	var currentUser models.User
	result := database.DB.Where("name = ?", cookie.Value).First(&currentUser)
	if result.Error != nil {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	commentTmpl, err := template.ParseFiles("public/comment.html")
	if err != nil {
		http.Error(w, "HTML not found", http.StatusSeeOther)
		return
	}
	commentTmpl.Execute(w, currentUser)
}

// POST /comment_action
func CommentAction(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		inputText := r.FormValue("comment_text_input")

		cookie, err := r.Cookie("user_session")
		if err != nil {
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
			return
		}

		var currentUser models.User
		userResult := database.DB.Where(&models.User{Name: cookie.Value}).First(&currentUser)
		if userResult.Error != nil {
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
			return
		}

		var commentCount int64
		database.DB.Model(&models.Comment{}).Where("user_id = ?", currentUser.ID).Count(&commentCount)

		if commentCount >= 5 {
			http.Error(w, "You can maximum write 5 Comments", http.StatusBadRequest)
			return
		}

		newComment := models.Comment{
			Text:   inputText,
			UserID: currentUser.ID,
		}

		result := database.DB.Create(&newComment)
		if result.Error != nil {
			http.Error(w, "Problem with Create Comment", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}