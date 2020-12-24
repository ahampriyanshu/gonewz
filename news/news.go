package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

/*Article : Structure to store the json object*/
type Article struct {
	Source struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}

/*FormatPublishedDate : Formatting published date*/
func (a *Article) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}

/*Results : Structure to store the results */
type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

/*Client : Structure to store http client */
type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

/*FetchHeadline : To fetch the headlines*/
func (c *Client) FetchHeadline(page string) (*Results, error) {
	endpoint := fmt.Sprintf("https://saurav.tech/NewsAPI/top-headlines/category/general/in.json")
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{}
	return res, json.Unmarshal(body, res)
}

/*FetchByCategory : Fectch for Category */
func (c *Client) FetchByCategory(category, page string) (*Results, error) {
	endpoint := fmt.Sprintf("https://saurav.tech/NewsAPI/top-headlines/category/%s/in.json", url.QueryEscape(category))
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{}
	return res, json.Unmarshal(body, res)
}

/*FetchBySearch : To fetch the queries*/
func (c *Client) FetchBySearch(query, page string) (*Results, error) {
	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(query), c.PageSize, page, c.key)
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{}
	return res, json.Unmarshal(body, res)
}

/*FetchBySource : To fetch by source*/
func (c *Client) FetchBySource(query, page string) (*Results, error) {
	endpoint := fmt.Sprintf("https://saurav.tech/NewsAPI/everything/%s.json", url.QueryEscape(query))
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{}
	return res, json.Unmarshal(body, res)
}

/*NewClient : <aking requests to the News api. */
func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}
