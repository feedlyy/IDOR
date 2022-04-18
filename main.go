package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

var schema = `
CREATE TABLE person (
	id SERIAL PRIMARY KEY,
    first_name text,
    last_name text,
    email text
);
CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	price NUMERIC(10,2) NOT NULL,
	date_time TIMESTAMP WITH TIME ZONE NOT NULL
);`

type Person struct {
	Id        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
}

type Order struct {
	Id       int    `db:"id" json:"id"`
	UserId   int    `db:"user_id" json:"user_id"`
	Name     string `db:"name" json:"name"`
	Price    string `db:"price" json:"price"`
	DateTime string `db:"date_time" json:"date_time"`
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		"localhost", 5432, "fadli", "nill", "local")
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	if os.Getenv("migrate") == "true" {
		db.MustExec(schema)
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
		tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
		tx.MustExec("INSERT INTO orders (user_id, name, price, date_time) VALUES ($1, $2, $3, $4)", 1, "iPhone", "100000", "2019-01-01 12:00:00")
		tx.MustExec("INSERT INTO orders (user_id, name, price, date_time) VALUES ($1, $2, $3, $4)", 1, "Macbook", "150000", "2019-01-01 12:00:00")
		tx.MustExec("INSERT INTO orders (user_id, name, price, date_time) VALUES ($1, $2, $3, $4)", 2, "Ipad", "150000", "2019-01-01 12:00:00")
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}

	router := httprouter.New()
	router.GET("/order/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		var (
			id        = params.ByName("id")
			order     Order
			userId, _ = strconv.Atoi(request.Header.Get("user-id"))
		)
		fmt.Println(userId)

		//err = db.Get(&order, "SELECT * FROM orders WHERE id = $1", id)
		err = db.Get(&order, "SELECT * FROM orders WHERE user_id = $1 and id = $2", userId, id)
		if err != nil {
			ResponseError(writer, http.StatusNotFound, err)
			return
		}

		ResponseJSON(writer, http.StatusOK, order)
	})
	log.Print("server started at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
