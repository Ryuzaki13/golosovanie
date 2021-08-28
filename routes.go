package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"time"
)

func index(c *gin.Context) {
	Branch, Phone, isStudent, err := get_user_info()
	if err != nil {
		c.Data(200, "text/html; charset=utf-8", []byte("Пользователь не авторизирован"))
	} else if Branch == "44b72cd889a44234a8a7a49750fedf1bc5a654d8b06257ec5866711be0be6286" || isStudent == false {
		voting(c, Phone, isStudent)
	} else {
		c.Data(200, "text/html; charset=utf-8", []byte("Только для студентов АиВТ"))
	}
	fmt.Println(c.Request.Host + c.FullPath())
}

func voting(c *gin.Context, user_id int, user_student bool) {
	var err error
	var vote vote_type

	vote.Date, err = time.Parse("2006-01-02", c.Param("date"))
	hist := true
	if err != nil {
		hist = false
		vote.Date = get_date()
	}

	err = connection.QueryRow(context.Background(), `
    SELECT voting_id, date, winner
    FROM history
    WHERE date = $1
    `, vote.Date.Format("2006-01-02")).Scan(&vote.Voting_id, &vote.Date, &vote.Winner)
	if err != nil {
		if err.Error() == string("no rows in result set") {
			err = connection.QueryRow(context.Background(), `
            INSERT INTO history (date)
            VALUES ($1);

			INSERT INTO votes (
				SELECT (SELECT voting_id FROM history WHERE date = $1 LIMIT 1),song_id,ARRAY[]::integer[],0
				FROM songs WHERE active = true
			);

			SELECT voting_id, date, winner
            FROM history
            WHERE date = $1;
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
    WHERE voting_id = $1 AND songs.song_id = votes.song_id
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

	// json1, err := json.MarshalIndent(vote, "", "    ")
	// fmt.Println("vote " + string(json1) + "\n")

	user_choice := 0
Loop:
	for _, e1 := range vote.Votes {
		for _, e2 := range e1.Votes {
			if user_id == e2 {
				user_choice = e1.Song_id
				break Loop
			}
		}
	}

	c.HTML(
		http.StatusOK,
		"index.html",
		struct {
			Vote_id      int
			Date         string
			Winner       int
			Votes        []votes_type
			User_id      int
			User_choice  int
			User_student bool
			Hist         bool
		}{
			vote.Voting_id,
			vote.Date.Format("2006-01-02"),
			vote.Winner,
			vote.Votes,
			user_id,
			user_choice,
			user_student,
			hist,
		},
	)
}

func choice(c *gin.Context) {
	var points int
	_, _, isStudent, _ := get_user_info()
	if isStudent == true {
		points = ves_s
	} else {
		points = ves_p
	}

	body, _ := ioutil.ReadAll(c.Request.Body)
	var res struct {
		User_id   int `json:"id"`
		Voting_id int `json:"voting_id"`
		Choice    int `json:"choice"`
	}
	err := json.Unmarshal(body, &res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Error1: %v\n", err)
	} else {
		_, err = connection.Exec(context.Background(), `
		UPDATE votes
		SET points = (
			SELECT (
				(SELECT points FROM votes WHERE ARRAY[$1::int] <@ votes AND voting_id = $2 LIMIT 1) - $3::int
			)
		), votes = (
			SELECT array(
				SELECT unnest(
					array(SELECT votes FROM votes WHERE ARRAY[$1::int] <@ votes AND voting_id = $2 LIMIT 1)
				)
				EXCEPT SELECT $1::int
			)
		)
		WHERE (ARRAY[$1::int] <@ votes) AND voting_id = $2;
		`, res.User_id, res.Voting_id, points)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Query Error2: %v\n", err)
		} else {
			_, err = connection.Exec(context.Background(), `
			UPDATE votes
			SET points = (
				SELECT (
					(SELECT points FROM votes WHERE voting_id = $2 AND song_id = $3 LIMIT 1) + $4::int
				)
			), votes = array_append(votes, $1)
			WHERE voting_id = $2 AND song_id = $3;
			`, res.User_id, res.Voting_id, res.Choice, points)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Query Error3: %v\n", err)
			}
		}
	}
	if err != nil {
		c.Data(200, "text/html; charset=utf-8", []byte(`0`))
	} else {
		c.Data(200, "text/html; charset=utf-8", []byte(`1`))
	}
}

func history(c *gin.Context) {
	date := get_date()

	rows, err := connection.Query(context.Background(), `
    SELECT history.date,songs.name,songs.url
	FROM history,songs
	WHERE songs.song_id = history.winner
	ORDER BY history.date DESC
    `)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Error1: %v\n", err)
	}
	defer rows.Close()

	if rows.Err() != nil {
		fmt.Fprintf(os.Stderr, "Query Error2: %v\n", err)
	}

	type hist []struct {
		Date string
		Name string
		Url  string
	}

	var h hist

	for rows.Next() {
		var v struct {
			Date time.Time
			Name string
			Url  string
		}
		err = rows.Scan(&v.Date, &v.Name, &v.Url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query Error5: %v\n", err)
		}
		h = append(h, struct {
			Date string
			Name string
			Url  string
		}{
			v.Date.Format("2006-01-02"), v.Name, v.Url,
		})
	}

	c.HTML(
		http.StatusOK,
		"history.html",
		struct {
			Date    string
			History hist
		}{
			date.Format("2006-01-02"),
			h,
		},
	)
}

func get_songs(c *gin.Context) {
	_, _, isStudent, _ := get_user_info()
	if isStudent == true {
		c.Data(200, "text/html; charset=utf-8", []byte("Только для преподавателей"))
	} else {
		rows, err := connection.Query(context.Background(), `
    	SELECT name,url,active FROM songs OFFSET 1
    	`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query Error3: %v\n", err)
		}
		defer rows.Close()

		if rows.Err() != nil {
			fmt.Fprintf(os.Stderr, "Query Error4: %v\n", err)
		}

		type song struct {
			Name   string
			Url    string
			Active bool
		}

		var arr []song

		for rows.Next() {
			var s song
			err = rows.Scan(&s.Name, &s.Url, &s.Active)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Query Error5: %v\n", err)
			}
			arr = append(arr, struct {
				Name   string
				Url    string
				Active bool
			}{
				s.Name, s.Url, s.Active,
			})
		}

		c.HTML(
			http.StatusOK,
			"songs.html",
			struct {
				Songs []song
			}{
				arr,
			},
		)
	}
}

func edit_songs(c *gin.Context) {
	_, _, isStudent, _ := get_user_info()
	if isStudent == true {
		c.Data(200, "text/html; charset=utf-8", []byte("Только для преподавателей"))
	} else {
		body, _ := ioutil.ReadAll(c.Request.Body)
		// fmt.Println(string(body))
		var res [][]interface{}
		err := json.Unmarshal(body, &res)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query Error1: %v\n", err)
		} else {
			sql := "DELETE FROM songs WHERE song_id != 0;"
			for _, v := range res {
				sql = sql + fmt.Sprintf(`
				INSERT INTO "songs" ("name","url","active")
				VALUES ('%v','%v',%v);
				`, v[0], v[1], v[2])
			}
			// fmt.Println(sql)
			_, err = connection.Exec(context.Background(), sql, res)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Query Error2: %v\n", err)
			}
		}
		if err != nil {
			c.Data(200, "text/html; charset=utf-8", []byte(`0`))
		} else {
			c.Data(200, "text/html; charset=utf-8", []byte(`1`))
		}
	}
}
