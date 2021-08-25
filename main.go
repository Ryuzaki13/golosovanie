package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ves_s int = 1
var ves_p int = 10

type votes_type struct {
	Song_id int    `json:"Song_id"`
	Votes   []int  `json:"Votes"`
	Name    string `json:"Name"`
	Url     string `json:"Url"`
}

type vote_type struct {
	Voting_id int          `json:"Voting_id"`
	Date      time.Time    `json:"Date"`
	Winner    int          `json:"Winner"`
	Votes     []votes_type `json:"Votes"`
}

func main() {

	connection, err := pgxpool.Connect(context.Background(), "postgres://postgres:1@localhost:5432/postgres?statement_cache_capacity=0&pool_max_conns=1000")
	// min_read_buffer_size, servicefile, statement_cache_capacity, statement_cache_mode, prefer_simple_protocol, pool_min_conns, pool_max_conns, pool_max_conn_lifetime, pool_max_conn_idle_time, pool_health_check_period
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer connection.Close()

	file, _ := os.Open("sql.sql")
	str, _ := ioutil.ReadAll(file)
	_, err = connection.Exec(context.Background(), string(str))

	// connection.Exec() - для выполнения запроса, который не возвращает набор результатов
	// connection.Query() - получить ряды
	// connection.QueryRow() - получить 1 ряд
	// connection.QueryFunc() - execute a callback function for every row

	//

	//

	//

	// создание роутера
	var router_setup = func() *gin.Engine {

		var router = gin.Default()

		router.Static("/static/", "./static")

		var formatAsDate = func(t time.Time) string {
			year, month, day := t.Date()
			return fmt.Sprintf("%d%02d/%02d", year, month, day)
		}
		var arr_int_len = func(v []int) int {
			return len(v)
		}
		router.SetFuncMap(template.FuncMap{
			"formatAsDate": formatAsDate,
			"arr_int_len":  arr_int_len,
		})
		router.LoadHTMLGlob("templates/*")

		// данные аккаунта
		router.POST("/api/user/select-info", func(c *gin.Context) {
			body := []byte(`{
			"Data": {
				"Branch": "44b72cd889a44234a8a7a49750fedf1bc5a654d8b06257ec5866711be0be6286",
				"Phone": "12345",
				"isStudent": false
			},
			"Error": null
		}`)
			c.Data(200, "application/json; charset=utf-8", body)
		})

		// начало
		router.GET("/", func(c *gin.Context) {
			c.Data(200, "text/html; charset=utf-8", []byte(`<script src="./static/1.js"></script>`))
			fmt.Println(c.Request.Host + c.FullPath())
		})

		// изменение списка песен
		router.GET("/edit-songs", func(c *gin.Context) {
		})
		router.POST("/edit-songs", func(c *gin.Context) {
		})

		//выбор песни в текущем голосовании
		router.POST("/choice", func(c *gin.Context) {
		})

		// история голосований
		// router.POST("/history", func(c *gin.Context) {
		// })

		// текущее голосование
		router.POST("/vote", func(c *gin.Context) {
			req_body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				fmt.Println("req_body error\n", err)
			}
			a := strings.Split(string(req_body[:]), "-")
			user_id, _ := strconv.Atoi(a[0])
			user_student, _ := strconv.ParseBool(a[1])

			var vote vote_type
			// date := time.Now().AddDate(0, 0, 7-int(time.Now().Weekday()))
			vote.Date, _ = time.Parse("2006-01-02", "2021-01-03")

			err = connection.QueryRow(context.Background(), `SELECT voting_id, date, winner FROM history WHERE date = $1`, vote.Date.Format("2006-01-02")).Scan(&vote.Voting_id, &vote.Date, &vote.Winner)
			if err != nil {
				if fmt.Sprint(err) == string("no rows in result set") {
					fmt.Println("+", vote.Date)
					connection.Exec(context.Background(), `INSERT INTO "history" ("date") VALUES ($1);`, vote.Date.Format("2006-01-02"))
					err = connection.QueryRow(context.Background(), `SELECT voting_id, date, winner FROM history WHERE date = $1`, vote.Date.Format("2006-01-02")).Scan(&vote.Voting_id, &vote.Date, &vote.Winner)
					if err != nil {
						fmt.Fprintf(os.Stderr, "QueryRow Error1: %v\n", err)
					}
					//!!!
				} else {
					fmt.Fprintf(os.Stderr, "QueryRow Error2: %v\n", err)
				}
			}

			rows, err := connection.Query(context.Background(), "SELECT song_id, votes FROM votes WHERE voting_id = $1", vote.Voting_id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Query Error3: %v\n", err)
			}
			defer rows.Close()

			if rows.Err() != nil {
				fmt.Fprintf(os.Stderr, "Query Error4: %v\n", err)
			}
			for rows.Next() {
				var Votes votes_type
				err = rows.Scan(&Votes.Song_id, &Votes.Votes)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Query Error5: %v\n", err)
				}
				vote.Votes = append(vote.Votes, Votes)
			}

			for i := 0; i < len(vote.Votes); i++ {
				connection.QueryRow(context.Background(), `SELECT name, url FROM songs WHERE song_id = $1`, vote.Votes[i].Song_id).Scan(&vote.Votes[i].Name, &vote.Votes[i].Url)
			}

			json1, err := json.MarshalIndent(vote, "", "    ")
			fmt.Println("vote " + string(json1) + "\n")

			user_choice := 0
		Loop:
			for i1, e1 := range vote.Votes {
				for _, e2 := range e1.Votes {
					if user_id == e2 {
						user_choice = vote.Votes[i1].Song_id
						break Loop
					}
				}
			}

			var data = struct {
				Date        string
				User_choice int
				User_id     int
				User_student bool
				Winner      int
				Votes       []votes_type
			}{
				vote.Date.Format("2006-01-02"),
				user_choice,
				user_id,
				user_student,
				vote.Winner,
				vote.Votes,
			}

			c.HTML(
				http.StatusOK,
				"index.html",
				data,
			)
		})

		return router
	}

	router_setup().Run()

	//

	//

	//

}
