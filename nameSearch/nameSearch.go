package nameSearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func NameSearch(search string) Data {
	type Variables struct {
		Search  string `json:"search"`
		Page    int    `json:"page"`
		PerPage int    `json:"perPage"`
	}

	jsonData := struct {
		Query     string    `json:"query"`
		Variables Variables `json:"variables"`
	}{
		Query: `query ($id: Int, $page: Int, $perPage: Int, $search: String) {
			Page (page: $page, perPage: $perPage) {
				pageInfo {
					total
					currentPage
					lastPage
					hasNextPage
					perPage
				}
				media (id: $id, search: $search) {
					id
					title {
						romaji
					}
				}
			}
		}`,
		Variables: Variables{
			Search:  search,
			Page:    1,
			PerPage: 10,
		},
	}
	test := http.Client{}
	jsonValue, _ := json.Marshal(jsonData)
	req, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := test.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()
	var data Data
	json.NewDecoder(resp.Body).Decode(&data)
	return data
}

type Data struct {
	Data struct {
		Page struct {
			PageInfo struct {
				Total       int  `json:"total"`
				CurrentPage int  `json:"currentPage"`
				LastPage    int  `json:"lastPage"`
				HasNextPage bool `json:"hasNextPage"`
				PerPage     int  `json:"perPage"`
			} `json:"pageinfo"`
			Media []struct {
				Id    int `json:"id"`
				Title struct {
					Romaji string `json:"romaji"`
				} `json:"title"`
			} `json:"media"`
		} `json:"page"`
	} `json:"data"`
	Errors []struct {
		Message   string `json:"message"`
		Status    int    `json:"status"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
	}
}

// const doc = `<!DOCTYPE html>
// <html>
//     <head>
// 	<title>Name Searcher</title>
//     </head>
//     <body>
// 	<p> Page Info:</p>
// 	<ul>{{with .Data.Page.PageInfo}}
// 		<li>Total		:	{{.Total}} </li>
// 		<li>CurrentPage	:	{{.CurrentPage}} </li>
// 		<li>LastPage	:	{{.LastPage}} </li>
// 		<li>HasNextPage	:	{{.HasNextPage}} </li>
// 		<li>PerPage		:	{{.PerPage}} </li>
// 		{{end}}
// 	</ul>
// 	<p>Data: </p>
// 	<ul> {{range .Data.Page.Media}}
// 		<li>ID		:	{{.Id}}  |  Title	:	{{.Title.Romaji}} </li>
// 		{{end}}
// 	</ul>
//     </body>
// </html>
// `
