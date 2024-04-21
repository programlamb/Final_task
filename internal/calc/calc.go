package calc

import (
	"go/token"
	"go/types"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"Final_task/internal/authorization"
	"Final_task/internal/db"
)

// map хранящий настройки времени выполнения операций в секундах
var settings = map[string]int{"+": 1, "-": 1, "*": 1, "/": 1}

// Handler главной страницы
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if authorization.GetActiveUserID() == 0 {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// объект с HTML кодом главной страницы
	tmplIndex := template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))

	// проверяем должны ли выражения быть посчитаны к моменту загрузки страницы
	// если да, то считаем их
	exps := updateStatus(db.GetExpressions(authorization.GetActiveUserID()))
	// загружаем страницу по HTML шаблону и передаём выражения, чтобы отобразить их на странице
	tmplIndex.ExecuteTemplate(w, "index", exps)
}

// Handler, который обрабатывает форму с выражением
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	expression := r.FormValue("exp")

	// рассчитывем время, когда выражение посчитается путём сложения времени выполнения всех операций
	createTime := time.Now()
	calcTime := createTime
	for i := 0; i < len(expression); i++ {
		if v, ok := settings[string(expression[i])]; ok {
			calcTime = calcTime.Add(time.Second * time.Duration(v))
		}
	}

	exp := db.Expression{Exp: expression, CalcTime: calcTime, CreateTime: createTime}
	db.AddExpression(authorization.GetActiveUserID(), exp)
	// переводим пользователя обратно на главную страницу
	http.Redirect(w, r, "http://localhost:8080/", http.StatusPermanentRedirect)
}

// Handler страницы настроек
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if authorization.GetActiveUserID() == 0 {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// объект с HTML кодом страницы настроек
	tmplSettings := template.Must(template.ParseFiles("templates/settings.html", "templates/base.html"))

	// загружаем страницу по HTML шаблону и передаём настройки, чтобы отобразить их на странице в строках ввода
	// но там всегда будут стоять значения 1 секунда, потому что я это не реализовал
	tmplSettings.ExecuteTemplate(w, "settings", settings)
}

// Handler, который обрабатывает форму с настройками
func SettingsSaver(w http.ResponseWriter, r *http.Request) {
	// Было бы круто сделать это через пост реквест, но я делаю это через параметры в адресе
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	for key := range settings {
		value, err := strconv.Atoi(params.Get(key))
		if err == nil {
			settings[key] = value
		}
	}

	// переводим пользователя обратно на главную с настройками
	http.Redirect(w, r, "http://localhost:8080/settings", http.StatusPermanentRedirect)
}

// функция, которая обновляет статус и решает выражения
func updateStatus(e []db.Expression) []db.Expression {
	// проверяем всё выраженя, и если прошёл момент их вычисления, вычисляем выражения
	for i := 0; i < len(e); i++ {
		if time.Since(e[i].CalcTime) >= 0 {
			// отмечаем, что выражение посчитано
			e[i].Status = ""

			// я не уверен, что это за fileSet, но он нужен для разбиения примера на операции, числа и скобки
			fs := token.NewFileSet()
			// разбиваем пример на операции, числа и скобки
			tv, err := types.Eval(fs, nil, token.NoPos, e[i].Exp)
			if err != nil {
				e[i].Result = ""
				e[i].Status = "expression parsing error"
			} else {
				// считаем пример и записываем его значение
				e[i].Result = tv.Value.String()
			}
		} else {
			e[i].Status = "expression will be calculated soon"
		}
	}
	return e
}
