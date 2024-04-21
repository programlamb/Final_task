package main

import (
	"net/http"

	"Final_task/internal/authorization"
	"Final_task/internal/calc"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // подключение CSS файла к серверу для визуала страницы

	// подключение hendler-ов к серверу по заданым путям
	mux.HandleFunc("/register", authorization.RegisterHandler)
	mux.HandleFunc("/save_user", authorization.SaveUserHandler)
	mux.HandleFunc("/login", authorization.LoginHandler)
	mux.HandleFunc("/loginAsUser", authorization.LoginAsUserHandler)
	mux.HandleFunc("/exit", authorization.Exit)
	mux.HandleFunc("/settings", calc.SettingsHandler)
	mux.HandleFunc("/setSettings", calc.SettingsSaver)
	mux.HandleFunc("/calc", calc.CalcHandler)
	mux.HandleFunc("/", calc.IndexHandler)

	// запускаем сервер по адресу http://localhost:8080/
	http.ListenAndServe(":8080", mux)
}
