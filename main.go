package main

import (
	"fmt"
	"net/http"
	"os"
	"shadow/anilists/test/addNames"
	"shadow/anilists/test/auth"
	"shadow/anilists/test/nameSearch"
	"text/template"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}
func indexPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/index.html", indexPage)
	// http.HandleFunc("/id", getID)
	http.HandleFunc("/name", getName)
	http.HandleFunc("/addname", addName)
	http.HandleFunc("/auth", autho)
	http.ListenAndServe(":8000", nil)
}

//	func getID(w http.ResponseWriter, r *http.Request) {
//		err := r.ParseForm()
//		if err != nil {
//			fmt.Println(err)
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		id, _ := strconv.Atoi(r.Form.Get("id"))
//		idSearch.IdSearch(id, w)
//	}
func getName(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	name := r.Form.Get("name")
	nameSearch.NameSearch(name)
}
func addName(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	query := r.Form.Get("addname")
	addNames.AddNames(query)
	fo, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	_, err = fo.WriteString(query + "\n")
	if err != nil {
		panic(err)
	}

}
func autho(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userName := r.Form.Get("userName")
	code := r.Form.Get("auth")
	auth.Auth(userName, code, w)
}

// func serveTemplate(w http.ResponseWriter, r *http.Request) {
// 	lp := filepath.Join("templates", "layout.html")
// 	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

// 	// Return a 404 if the template doesn't exist
// 	info, err := os.Stat(fp)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			http.NotFound(w, r)
// 			return
// 		}
// 	}

// 	// Return a 404 if the request is for a directory
// 	if info.IsDir() {
// 		http.NotFound(w, r)
// 		return
// 	}

// 	tmpl, err := templates.ParseFiles(lp, fp)
// 	if err != nil {
// 		// Log the detailed error
// 		log.Print(err.Error())
// 		// Return a generic "Internal Server Error" message
// 		http.Error(w, http.StatusText(500), 500)
// 		return
// 	}
// 	err = tmpl.ExecuteTemplate(w, "layout", nil)
// 	if err != nil {
// 		log.Print(err.Error())
// 		http.Error(w, http.StatusText(500), 500)
// 	}
// }

// func landingPage(w http.ResponseWriter, r *http.Request) {
// 	templates.ExecuteTemplate(w, "home.html", nil)
// }

// func newRouter() *mux.Router {
// 	r := mux.NewRouter()

// 	staticFileDirectory := http.Dir("./assets/")
// 	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))

// 	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET", "POSt")
// 	r.HandleFunc("/id", getID).Methods("GET", "POST")
// 	r.HandleFunc("/name", getName).Methods("GET", "POST")
// 	r.HandleFunc("/add", addName).Methods("GET", "POST")
// 	r.HandleFunc("/auth", autho).Methods("GET", "POST")
// 	return r
// }
