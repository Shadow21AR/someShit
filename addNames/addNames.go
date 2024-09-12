package addNames

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"shadow/anilists/test/nameSearch"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sado"
	password = "sado@123"
	dbname   = "sado"
)

var markReadQuery = `
mutation ($mediaId: Int) {
	SaveMediaListEntry (mediaId: $mediaId, status: CURRENT) {
		id
		status
	}
}
`

func AddNames(query string) {
	var db *sql.DB
	var err error
	var token string
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	row := db.QueryRow(`SELECT token FROM token WHERE username='shadow21a';`)
	err = row.Scan(&token)
	if err != nil {
		panic(err)
	}
	search := nameSearch.NameSearch(query)
	fmt.Printf("%#v", search)

	body := map[string]interface{}{
		"query": markReadQuery,
		"variables": map[string]interface{}{
			"mediaId": search.Data.Page.Media[0].Id,
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		"https://graphql.anilist.co",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	fmt.Println("Sending request to Anilist: ")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	var response struct {
		Data struct {
			SaveMediaListEntry struct {
				ID int `json:"ID"`
			} `json:"SaveMediaListEntry"`
		} `json:"data"`
		Errors []struct {
			Message   string `json:"message"`
			Status    int    `json:"status"`
			Locations []struct {
				Line   int `json:"line"`
				Column int `json:"column"`
			} `json:"locations"`
		} `json:"errors"`
	}
	json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println("Response: ", response.Data.SaveMediaListEntry.ID)
	fmt.Println("Error: ", response.Errors[0].Message)
}
