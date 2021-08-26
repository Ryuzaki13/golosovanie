package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"

	"time"
)

//

// получение данных пользователя
func get_user_info() (string, int, bool, error) {
	Branch := "44b72cd889a44234a8a7a49750fedf1bc5a654d8b06257ec5866711be0be6286"
	Phone := 12345
	isStudent := false
	var Error error
	return Branch, Phone, isStudent, Error
}

// подключение к бд
func connect_to_db(url string) *pgxpool.Pool {
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

var connection = connect_to_db("postgres://postgres:1@localhost:5432/postgres?statement_cache_capacity=0&pool_max_conns=1000")

// min_read_buffer_size, servicefile, statement_cache_capacity, statement_cache_mode, prefer_simple_protocol, pool_min_conns, pool_max_conns, pool_max_conn_lifetime, pool_max_conn_idle_time, pool_health_check_period

// connection.Exec() - для выполнения запроса, который не возвращает набор результатов
// connection.Query() - получить ряды
// connection.QueryRow() - получить 1 ряд
// connection.QueryFunc() - execute a callback function for every row

func index(c *gin.Context) {
	Branch, Phone, isStudent, err := get_user_info()
	if err != nil {
		c.Data(200, "text/html; charset=utf-8", []byte("Пользователь не авторизирован"))
	} else if Branch == "44b72cd889a44234a8a7a49750fedf1bc5a654d8b06257ec5866711be0be6286" || isStudent == false {
		current_vote(c, Phone, isStudent)
	} else {
		c.Data(200, "text/html; charset=utf-8", []byte("Только для студентов АиВТ"))
	}
	// fmt.Println(c.Request.Host + c.FullPath())
}

func current_vote(c *gin.Context, user_id int, user_student bool) {

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

	var connection = connect_to_db("postgres://postgres:1@localhost:5432/postgres?statement_cache_capacity=0&pool_max_conns=1000")
	defer connection.Close()

	var vote vote_type
	// date := time.Now().AddDate(0, 0, 7-int(time.Now().Weekday()))
	vote.Date, _ = time.Parse("2006-01-02", "2021-01-03")

	err := connection.QueryRow(context.Background(), `
    SELECT voting_id, date, winner
    FROM history
    WHERE date = $1
    `, vote.Date.Format("2006-01-02")).Scan(&vote.Voting_id, &vote.Date, &vote.Winner)
	if err != nil {
		if fmt.Sprint(err) == string("no rows in result set") {
			connection.Exec(context.Background(), `
            INSERT INTO "history" ("date")
            VALUES ($1)
            `, vote.Date.Format("2006-01-02"))
			err = connection.QueryRow(context.Background(), `
            SELECT voting_id, date, winner
            FROM history
            WHERE date = $1
            `, vote.Date.Format("2006-01-02")).Scan(&vote.Voting_id, &vote.Date, &vote.Winner)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow Error1: %v\n", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "QueryRow Error2: %v\n", err)
		}
	}

	rows, err := connection.Query(context.Background(), `
    SELECT votes.song_id, votes.votes, votes.points, songs.name, songs.url
    FROM votes, songs
    WHERE voting_id = $1 AND votes.song_id = songs.song_id
    `, vote.Voting_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Error3: %v\n", err)
	}
	defer rows.Close()

	if rows.Err() != nil {
		fmt.Fprintf(os.Stderr, "Query Error4: %v\n", err)
	}

	for rows.Next() {
		var Votes votes_type
		err = rows.Scan(&Votes.Song_id, &Votes.Votes, &Votes.Points, &Votes.Name, &Votes.Url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query Error5: %v\n", err)
		}
		vote.Votes = append(vote.Votes, Votes)
	}

	for i := 0; i < len(vote.Votes); i++ {
		connection.QueryRow(context.Background(), `
        SELECT name, url
        FROM songs
        WHERE song_id = $1
        `, vote.Votes[i].Song_id).Scan(&vote.Votes[i].Name, &vote.Votes[i].Url)
	}

	json1, err := json.MarshalIndent(vote, "", "    ")
	fmt.Println("vote " + string(json1) + "\n")

	user_choice := 0
Loop:
	for i := range vote.Votes {
		for _, e := range vote.Votes[i].Votes {
			if user_id == e {
				user_choice = vote.Votes[i].Song_id
				break Loop
			}
		}
	}

	var data = struct {
		Vote_id      int
		Date         string
		Winner       int
		Votes        []votes_type
		User_id      int
		User_choice  int
		User_student bool
	}{
		vote.Voting_id,
		vote.Date.Format("2006-01-02"),
		vote.Winner,
		vote.Votes,
		user_id,
		user_choice,
		user_student,
	}

	c.HTML(
		http.StatusOK,
		"index.html",
		data,
	)
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

	// начало
	router.GET("/", index)

	// получение списка песен
	router.GET("/edit-songs", func(c *gin.Context) {
	})

	// изменение списка песен
	router.POST("/edit-songs", func(c *gin.Context) {
	})

	// выбор песни в текущем голосовании
	router.POST("/choice", func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		var res struct {
			User_id   int `json:"id"`
			Voting_id int `json:"voting_id"`
			Choice    int `json:"choice"`
		}
		err := json.Unmarshal(body, &res)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query Error6: %v\n", err)
		} else {
			_, err = connection.Exec(context.Background(), `
            UPDATE votes SET votes = (
                SELECT array(
                    SELECT unnest(
                        array(SELECT votes FROM votes WHERE (ARRAY[$1] <@ votes) AND voting_id = $2)
                    )
                    EXCEPT SELECT $1
                )
            )
            WHERE (ARRAY[$1] <@ votes) AND voting_id = $2;
            `, res.User_id, res.Voting_id)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Query Error7: %v\n", err)
			} else {
				_, err = connection.Exec(context.Background(), `
            UPDATE votes
            SET votes = array_append(votes, $1)
            WHERE voting_id = $2 AND song_id = $3
            `, res.User_id, res.Choice)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Query Error8: %v\n", err)
				}
			}
		}
	})

	// история голосований
	router.GET("/history", func(c *gin.Context) {
	})

	return router
}
