package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ves_s int = 1
var ves_p int = 10

func main() {

	// создание роутера
	var router_setup = func() *gin.Engine {
		var router = gin.Default()

		router.Static("/static/", "./static")
		router.LoadHTMLGlob("templates/*")

		// данные аккаунта
		router.POST("/api/user/select-info", func(c *gin.Context) {
			body := []byte(`{
				"Data": {
					"Branch": "44b72cd889a44234a8a7a49750fedf1bc5a654d8b06257ec5866711be0be6286",
					"Phone": "12345"
				},
				"Error": null
			}`)
			c.Data(http.StatusOK, "application/json; charset=utf-8", body)
		})

		// главная
		router.GET("/", func(c *gin.Context) {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<script src="./static/1.js"></script>`))
			// c.HTML(http.StatusOK, "index.html", gin.H{})
		})

		// текущее голосование
		router.POST("/current", func(c *gin.Context) {
			req_body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			id := string(req_body[:])

			// date := time.Now().AddDate(0, 0, 7-int(time.Now().Weekday())).Format("2006-01-02")
			date := "2021-01-17"

			fmt.Println("---id---: " + id)
			fmt.Println("---date---: " + date)

			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(`qwe`))
		})

		return router
	}

	if false {
		router_setup().Run()
	}

	connection, err := pgxpool.Connect(context.Background(), "postgres://postgres:1@localhost:5432/postgres?statement_cache_capacity=0&pool_max_conns=999")
	// min_read_buffer_size
	// servicefile
	// statement_cache_capacity
	// statement_cache_mode
	// prefer_simple_protocol
	// pool_min_conns
	// pool_max_conns
	// pool_max_conn_lifetime
	// pool_max_conn_idle_time
	// pool_health_check_period
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer connection.Close()

	// connection.Exec() //для выполнения запроса, который не возвращает набор результатов
	// connection.Query() //получить ряды
	// connection.QueryRow() //получить 1 ряд
	// connection.QueryFunc() //execute a callback function for every row

	var resp []int
	err = connection.QueryRow(context.Background(), `select array(select id from files);`).Scan(&resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow Error: %v\n", err)
	}
	defer connection.Close()
	fmt.Println(resp)


	rows, err := connection.Query(context.Background(), "select generate_series(1,$1)", 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Error: %v\n", err)
	}
	defer rows.Close()

	if rows.Err() != nil {
		fmt.Fprintf(os.Stderr, "Query Error: %v\n", err)
	}
	for rows.Next() {
		var n int32
		err = rows.Scan(&n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query Error: %v\n", err)
		}
	}

}
