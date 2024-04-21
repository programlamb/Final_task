package authorization

import (
	"fmt"
	"html/template"
	"net/http"

	"Final_task/internal/db"
)

var userID int64 = 0

// Handler страницы регистрации
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// объект с HTML кодом
	tmplRegister := template.Must(template.ParseFiles("templates/register.html", "templates/base.html"))

	tmplRegister.ExecuteTemplate(w, "register", nil)
}

// функция, обрабатывающая форму регистрации
func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")

	if password == r.FormValue("password2") {
		user := db.User{Email: r.FormValue("email"), Name: r.FormValue("userName"), Password: password}
		userID = db.AddUser(user)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// объект с HTML кодом
	tmplLogin := template.Must(template.ParseFiles("templates/login.html", "templates/base.html"))

	tmplLogin.ExecuteTemplate(w, "login", nil)
}

func LoginAsUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	id, user, err := db.GetUser(email)
	if err == nil && user.Password == password {
		userID = id
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Print(err)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Exit(w http.ResponseWriter, r *http.Request) {
	userID = 0
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func GetActiveUserID() int64 {
	return userID
}
