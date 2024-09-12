package idSearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func IdSearch(id string) Data {
	type Variables struct {
		Id string `json:"id"`
	}

	jsonData := struct {
		Query     string    `json:"query"`
		Variables Variables `json:"variables"`
	}{
		Query: `query ($id: Int) {
				Media(id: $id, type: MANGA) {
					id
					title {
						romaji
						english
						native
					}
					genres
					description
				}
			}`,
		Variables: Variables{
			Id: id,
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
	// templates := template.New("template")
	// templates.New("doc").Parse(doc)
	// templates.Lookup("doc").Execute(w, data.Data.Media)
}

type Data struct {
	Data struct {
		Media struct {
			ID          int               `json:"id"`
			Title       map[string]string `json:"title"`
			Genres      []string          `json:"genres"`
			Description string            `json:"description"`
		} `json:"Media"`
	} `json:"data"`
}

// const doc = `<!DOCTYPE html>
// <html>
//     <head>
//     </head>
//     <body>
//         <h3>ID : {{.ID}}</h3>
// 		<ul>
// 		{{range .Title}}
// 		<li>{{.}}</li>
// 		{{end}}
// 		</ul>
// 		<h4>{{.Genres}}</h4>
// 		<h4>{{.Description}}</h4>
//     </body>
// </html>`
