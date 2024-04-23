package authorization

import (
	"html/template"
	"net/http"

	"Final_task/internal/db"
)

// Переменная хранящая номер авторизованного пользователя. Плохое решение, потому что эта перемнная одна на всех одновременных пользователей.
var userID int64 = 0

// Handler страницы регистрации
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// объект с HTML кодом
	tmplRegister := template.Must(template.ParseFiles("templates/register.html", "templates/base.html"))

	// загружаем страницу по HTML шаблону
	tmplRegister.ExecuteTemplate(w, "register", nil)
}

// функция, обрабатывающая форму регистрации
func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	// Проверку данных формы можно было бы усложнить, но пока так
	password := r.FormValue("password")

	if password == r.FormValue("password2") {
		// Сохраняем пользователя в базе данных
		user := db.User{Email: r.FormValue("email"), Name: r.FormValue("userName"), Password: password}
		userID = db.AddUser(user)
	}

	// загружаем страницу по HTML шаблону
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handler страницы авторизации
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// объект с HTML кодом
	tmplLogin := template.Must(template.ParseFiles("templates/login.html", "templates/base.html"))

	// загружаем страницу по HTML шаблону
	tmplLogin.ExecuteTemplate(w, "login", nil)
}

// функция, обрабатывающая форму авторизации
func LoginAsUserHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем данные пользователя
	email := r.FormValue("email")
	password := r.FormValue("password")

	id, user, err := db.GetUser(email)
	if err == nil && user.Password == password {
		userID = id
		// переводим пользователя обратно на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// переводим пользователя обратно на страницу авторизации
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// функция, обрабатывающая выход из профиля
func Exit(w http.ResponseWriter, r *http.Request) {
	userID = 0
	// переводим пользователя на страницу авторизации
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// функция, возвращающая ID авторизировнного пользователя
func GetActiveUserID() int64 {
	return userID
}
