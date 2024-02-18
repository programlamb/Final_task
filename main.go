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

type Expression struct {
	Exp        string
	Result     string
	Status     string
	CreateTime time.Time
	CalcTime   time.Time
}

var setings = map[string]int{"+": 1, "-": 1, "*": 1, "/": 1}

var exps = []Expression{}

var tplIndex = template.Must(template.ParseFiles("tamplates/index.html"))

var tplSattings = template.Must(template.ParseFiles("tamplates/settings.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	exps = updateStatus(exps)
	tplIndex.Execute(w, exps)
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	expression := params.Get("exp")

	createTime := time.Now()
	calcTime := createTime
	for i := 0; i < len(expression); i++ {
		if v, ok := setings[string(expression[i])]; ok {
			calcTime = calcTime.Add(time.Second * time.Duration(v))
		}
	}

	exps = append(exps, Expression{Exp: expression, Result: "?", Status: "expression will be calculated soon", CalcTime: calcTime, CreateTime: createTime})
	http.Redirect(w, r, "http://localhost:8080/", http.StatusPermanentRedirect)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	tplSattings.Execute(w, setings)
}

func settingsSaver(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, "http://localhost:8080/settings", http.StatusPermanentRedirect)
}

func updateStatus(e []Expression) []Expression {
	for i := 0; i < len(e); i++ {
		if time.Since(e[i].CalcTime) >= 0 {
			e[i].Status = ""

			fs := token.NewFileSet()
			tv, err := types.Eval(fs, nil, token.NoPos, e[i].Exp)
			if err != nil {
				e[i].Result = ""
				e[i].Status = "expression parsing error"
			} else {
				e[i].Result = tv.Value.String()
			}
		}
	}
	return e
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/settings", settingsHandler)
	mux.HandleFunc("/setSettings", settingsSaver)
	mux.HandleFunc("/calc", calcHandler)
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", mux)
}
