package main

import (
	"go/token"
	"go/types"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// структура каждого выражения
type Expression struct {
	Exp        string // сам пример
	Result     string // результат
	Status     string // статус(выполнено ли выражение, есть ли внём ошобки)
	CreateTime time.Time // время начала счёта
	CalcTime   time.Time // время когда выражение становится посчитанным
}

// map хранящий настройки времени выполнения операций в секундах
var setings = map[string]int{"+": 1, "-": 1, "*": 1, "/": 1}

// список хранящий все выражения
var exps = []Expression{}

// объект с HTML кодом главной страницы
var tplIndex = template.Must(template.ParseFiles("tamplates/index.html"))

// объект с HTML кодом страницы настроек
var tplSattings = template.Must(template.ParseFiles("tamplates/settings.html"))

// Handler главной страницы
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// проверяем должны ли выражения быть посчитаны к моменту загрузки страницы
	// если да, то считаем их
	exps = updateStatus(exps)
	// загружаем страницу по HTML шаблону и передаём выражения, чтобы отобразить их на странице
	tplIndex.Execute(w, exps)
}

// Handler, который обрабатывает форму с выражением
func calcHandler(w http.ResponseWriter, r *http.Request) {
	// Было бы круто сделать это через пост реквест, но я делаю это через параметры в адресе
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	expression := params.Get("exp")

	// рассчитывем время, когда выражение посчитается путём сложения времени выполнения всех операций
	createTime := time.Now()
	calcTime := createTime
	for i := 0; i < len(expression); i++ {
		if v, ok := setings[string(expression[i])]; ok {
			calcTime = calcTime.Add(time.Second * time.Duration(v))
		}
	}

	exps = append(exps, Expression{Exp: expression, Result: "?", Status: "expression will be calculated soon", CalcTime: calcTime, CreateTime: createTime})
	// переводим пользователя обратно на главную страницу
	http.Redirect(w, r, "http://localhost:8080/", http.StatusPermanentRedirect)
}

// Handler страницы настроек
func settingsHandler(w http.ResponseWriter, r *http.Request) {
	// загружаем страницу по HTML шаблону и передаём настройки, чтобы отобразить их на странице в строках ввода
	// но там всегда будут стоять значения 1 секунда, потому что я это не реализовал
	tplSattings.Execute(w, setings)
}

// Handler, который обрабатывает форму с настройками
func settingsSaver(w http.ResponseWriter, r *http.Request) {
	// Было бы круто сделать это через пост реквест, но я делаю это через параметры в адресе
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	for key := range setings {
		value, err := strconv.Atoi(params.Get(key))
		if err == nil {
			setings[key] = value
		}
	}

	// переводим пользователя обратно на главную с настройками
	http.Redirect(w, r, "http://localhost:8080/settings", http.StatusPermanentRedirect)
}

// функция, которая обновляет статус и решает выражения
func updateStatus(e []Expression) []Expression {
	// проверяем всё выраженя, и если прошёл момент их вычисления, вычисляем выражения
	// косяк: оно считает даже посчитанные. исправлю.
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
		}
	}
	return e
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // подключение CSS файла к серверу для визуала страницы

	// подключение hendler-ов к серверу по заданым путям
	mux.HandleFunc("/settings", settingsHandler)
	mux.HandleFunc("/setSettings", settingsSaver)
	mux.HandleFunc("/calc", calcHandler)
	mux.HandleFunc("/", indexHandler)

	// запускаем сервер по адресу http://localhost:8080/
	http.ListenAndServe(":8080", mux)
}
