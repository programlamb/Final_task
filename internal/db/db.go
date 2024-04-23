package db

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// структура каждого выражения
type Expression struct {
	Exp        string    // сам пример
	Result     string    // результат
	Status     string    // посчитано или нет, есть ли ошибки
	CreateTime time.Time // время начала счёта
	CalcTime   time.Time // время когда выражение становится посчитанным
}

// структура пользователя
type User struct {
	Email    string
	Name     string
	Password string
}

// функция сохранения пользователя в базе данных (возврашает ID этого пользователя)
func AddUser(user User) int64 {
	// Открываем базу данных
	ctx, db := runDB()
	defer db.Close()

	var q = "INSERT INTO users (email, name, password) VALUES ($1, $2, $3)"

	// Выполнем SQL запрос
	result, err := db.ExecContext(ctx, q, user.Email, user.Name, user.Password)
	if err != nil {
		panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	return id
}

// функция полученя информации о пользователе по его email
func GetUser(email string) (int64, User, error) {
	// Открываем базу данных
	ctx, db := runDB()
	defer db.Close()

	var user User
	var id int64
	var q = "SELECT id, email, name, password FROM users WHERE email = $1"

	// Выполнем SQL запрос
	err := db.QueryRowContext(ctx, q, email).Scan(&id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return 0, User{}, err
	}

	return id, user, nil
}

// функция сохранения выражения с привязкой к пользователю
func AddExpression(userID int64, exp Expression) {
	// Открываем базу данных
	ctx, db := runDB()
	defer db.Close()

	var q = "INSERT INTO expressions (user, exp, create_time, calc_time) VALUES ($1, $2, $3, $4)"

	// Выполнем SQL запрос
	_, err := db.ExecContext(ctx, q, userID, exp.Exp, exp.CreateTime, exp.CalcTime)
	if err != nil {
		panic(err)
	}
}

// функция получения всех выражений пользователя
func GetExpressions(userID int64) []Expression {
	// Открываем базу данных
	ctx, db := runDB()
	defer db.Close()

	var expressions []Expression
	var q = "SELECT exp, create_time, calc_time FROM expressions WHERE user = $1"

	// Выполнем SQL запрос
	rows, err := db.QueryContext(ctx, q, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.Exp, &e.CreateTime, &e.CalcTime)
		if err != nil {
			panic(err)
		}
		expressions = append(expressions, e)
	}

	return expressions
}

// функция, открывающая базу данных
func runDB() (context.Context, *sql.DB) {
	ctx := context.TODO()
	// Получаем путь к папке проекта
	exFile, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root := strings.ReplaceAll(exFile, `\`, `/`)
	db, err := sql.Open("sqlite3", root+"/db/users.db")

	if err != nil {
		panic(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, db
}
