package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type votes_type struct {
	Song_id int    `json:"Song_id"`
	Votes   []int  `json:"Votes"`
	Points  int    `json:"Points"`
	Name    string `json:"Name"`
	Url     string `json:"Url"`
}

type vote_type struct {
	Voting_id int          `json:"Voting_id"`
	Date      time.Time    `json:"Date"`
	Winner    int          `json:"Winner"`
	Votes     []votes_type `json:"Votes"`
}

var ves_s = 1
var ves_p = 10
var connection = connect_to_db("postgres://postgres:1@localhost:5432/postgres?statement_cache_capacity=0&pool_max_conns=1000")

func main() {
	defer connection.Close()

	router_setup().Run()
}

// подключение к бд
func connect_to_db(url string) *pgxpool.Pool {
	// min_read_buffer_size, servicefile, statement_cache_capacity, statement_cache_mode, prefer_simple_protocol, pool_min_conns, pool_max_conns, pool_max_conn_lifetime, pool_max_conn_idle_time, pool_health_check_period

	// connection.Exec() - для выполнения запроса, который не возвращает набор результатов
	// connection.Query() - получить ряды
	// connection.QueryRow() - получить 1 ряд
	// connection.QueryFunc() - execute a callback function for every row

	connection, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	file, _ := os.Open("sql.sql")
	str, _ := ioutil.ReadAll(file)
	_, err = connection.Exec(context.Background(), string(str))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	return connection
}

// получение данных пользователя
func get_user_info() (string, int, bool, error) {
	Branch := "44b72cd889a44234a8a7a49750fedf1bc5a654d8b06257ec5866711be0be6286"
	Phone := 12345
	isStudent := false
	var Error error
	return Branch, Phone, isStudent, Error
}

// получение текущей даты
func get_date() time.Time {
	// t, _ := time.Now().AddDate(0, 0, 7-int(time.Now().Weekday()))
	t, _ := time.Parse("2006-01-02", "2021-01-17")
	return t
}

// создание роутера
func router_setup() *gin.Engine {

	var router = gin.Default()

	router.Static("/static/", "./static")

	// var arr_int_len1 = func(v []int) int {
	//     return len(v)
	// }
	// router.SetFuncMap(template.FuncMap{
	//     "arr_int_len": arr_int_len1,
	// })

	router.LoadHTMLGlob("templates/*")

	//

	//

	// текущее голосование
	router.GET("/", index)

	// выбор песни в текущем голосовании
	router.POST("/choice", choice)

	// история голосований
	router.GET("/history", history)

	// голосование из истории
	router.GET("/history/:date", index)

	// получение списка песен
	router.GET("/edit-songs", get_songs)

	// изменение списка песен
	router.POST("/edit-songs", edit_songs)

	return router
}
