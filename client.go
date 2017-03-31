package main

import (
	"bytes"
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
	url, err := c.createURL(path, query)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("client: error creating http request: %s", err.Error())
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", c.token))
	req.Header.Add("Accept", c.AcceptHeader)

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client: error making a http request: %s", err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("client: error reading api error response: %s", err.Error())
		}
		return fmt.Errorf("client: GitHub api error - status %d: %s", res.StatusCode, body)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return fmt.Errorf("client: error unmarshaling json response: %s", err.Error())
	}
	return nil
}

// Post performs a post request with the provided data marshaled into the json body
func (c *Client) Post(path string, data interface{}) error {
	url, err := c.createURL(path, nil)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		return fmt.Errorf("client: error marshaling request body: %s", err.Error())
	}
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("client: error creating http request: %s", err.Error())
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", c.token))
	req.Header.Add("Accept", c.AcceptHeader)

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %s", err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 204 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("client: error reading api error response: %s", err.Error())
		}
		return fmt.Errorf("client: GitHub api error - status %d: %s", res.StatusCode, body)
	}
	return nil
}

func (c *Client) createURL(path string, query map[string]string) (string, error) {
	parsedURL, err := url.Parse(BaseURL)
	if err != nil {
		return "", fmt.Errorf("client: error creating an url: %s", err.Error())
	}
	parsedURL.Path = path
	if len(query) > 0 {
		v := url.Values{}
		for param, value := range query {
			v.Add(param, value)
		}
		parsedURL.RawQuery = v.Encode()
	}
	return parsedURL.String(), nil
}
