package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ahampriyanshu/gonewz/news"
	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

/*Q structure of Q */
type Q struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
	Type       int
}

/*IsLastPage : checking if the current page is the last page */
func (s *Q) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

/*CurrentPage : fetching current page */
func (s *Q) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}

/*PreviousPage : fetching previous page */
func (s *Q) PreviousPage() int {
	return s.CurrentPage() - 1
}

func dataHandler(newsapi *news.Client, pageType int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()

		Query := params.Get("q")
		if Query == "" {
			Query = "Why am I so lonely"
		}
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		var results *news.Results

		switch pageType {
		case 1:
			results, err = newsapi.FetchHeadline(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case 2:
			results, err = newsapi.FetchBySearch(Query, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case 3:
			results, err = newsapi.FetchByCategory(Query, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case 4:
			results, err = newsapi.FetchBySource(Query, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			results, err = newsapi.FetchBySearch(Query, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		q := &Q{
			Query:      Query,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults / newsapi.PageSize))),
			Results:    results,
			Type:       pageType,
		}

		if ok := !q.IsLastPage(); ok {
			q.NextPage++
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

func sendSW(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("sw.js")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Write(data)
}

func sendManifest(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("manifest.json")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apiKey, 20)

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/sw.js", sendSW)
	mux.HandleFunc("/manifest.json", sendManifest)
	mux.HandleFunc("/", dataHandler(newsapi, 1))
	mux.HandleFunc("/search", dataHandler(newsapi, 2))
	mux.HandleFunc("/category", dataHandler(newsapi, 3))
	mux.HandleFunc("/source", dataHandler(newsapi, 4))
	http.ListenAndServe(":"+port, mux)
}
