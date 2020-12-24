package main

import (
	"bytes"
	"fmt"
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

/*Search structure of search */
type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
	Type       int
}

/*IsLastPage : checking if the current page is the last page */
func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

/*CurrentPage : fetching current page */
func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}

/*PreviousPage : fetching previous page */
func (s *Search) PreviousPage() int {
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

		searchQuery := params.Get("q")
		if searchQuery == "" {
			searchQuery = "Why am I so lonely"
		}
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		var results *news.Results

		if pageType == 1 {
			results, err = newsapi.FetchForIndex(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if pageType == 2 {
			results, err = newsapi.FetchEverything(searchQuery, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if pageType == 3 {
			results, err = newsapi.FetchCategory(searchQuery, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if pageType == 4 {
			results, err = newsapi.FetchBySource(searchQuery, page)
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

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults / newsapi.PageSize))),
			Results:    results,
			Type:       pageType,
		}

		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
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
	} else {
		fmt.Println("Service worker !")
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
