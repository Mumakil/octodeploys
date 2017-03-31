package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client can make requests to GitHub api
type Client struct {
	client       *http.Client
	token        string
	AcceptHeader string
}

// BaseURL GitHub API base url
var BaseURL = "https://api.github.com/"

// DefaultAcceptHeader which version of GitHub api to use
var DefaultAcceptHeader = "application/vnd.github.v3+json"

// NewClient constructs a new client
func NewClient(token string) *Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &Client{
		client:       c,
		token:        token,
		AcceptHeader: DefaultAcceptHeader,
	}
}

// Get performs a get request and unmarshals the data to response
func (c *Client) Get(path string, query map[string]string, response interface{}) error {
	uri, err := url.Parse(BaseURL)
	if err != nil {
		return err
	}
	uri.Path = path
	if len(query) > 0 {
		v := url.Values{}
		for param, value := range query {
			v.Add(param, value)
		}
		uri.RawQuery = v.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", c.token))
	req.Header.Add("Accept", DefaultAcceptHeader)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("client: error reading api error response: %s", err.Error())
		}
		return fmt.Errorf("client: GitHub api error - status %d: %s", res.StatusCode, body)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&response)
	return err
}
