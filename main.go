package main

var ves_s int = 1
var ves_p int = 10

//

func voting_end() {
	// `SELECT * FROM "votes" WHERE "voting_id" = '$1' ORDER BY "points" DESC LIMIT 1`
}

func main() {

	router_setup().Run()

}
