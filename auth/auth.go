package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sado"
	password = "sado@123"
	dbname   = "sado"
)

func Auth(userName string, code string, w http.ResponseWriter) {
	body := map[string]interface{}{
		"client_id":     9658,
		"client_secret": "HZ5UjEzvD3rMZ0zXUtEH4AW7maAuY2QDhwykLkbn",
		"grant_type":    "authorization_code",
		"redirect_uri":  "https://anilist.co/api/v2/oauth/pin",
		"code":          code,
	}
	jsonBody, err := json.Marshal(&body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sending login request to Anilist")
	req, err := http.NewRequest(http.MethodPost, "https://anilist.co/api/v2/oauth/token", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Request fail %s\n", err)
	}
	fmt.Println("Decodint anilist shit")
	defer resp.Body.Close()
	var token Token
	json.NewDecoder(resp.Body).Decode(&token)

	fmt.Println("Got TOken...db shit ...")
	var db *sql.DB
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
	insert := `insert into "token"("username", "token") values($1, $2)`
	_, err = db.Exec(insert, userName, token.AccessToken)
	if err != nil {
		panic(err)
	}

	fmt.Println("Added to db!")
}

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type Post struct {
	Grant_type    string `json:"grant_type"`
	Client_id     string `json:"client_id"`
	Client_secret string `json:"client_secret"`
	Redirect_uri  string `json:"redirect_uri"`
	Code          string `json:"code"`
}
