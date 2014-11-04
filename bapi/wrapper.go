package bapi

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "errors"
    "fmt"
)

var (
    ApiUrl = "https://api.bitcoinaverage.com/"
)

type ApiClient struct {
    url         string
}

type Currency struct {
    Currency    string      `json:"currency"`
    Country     string      `json:"country"`
}

func New() *ApiClient {
    return NewWithOptions(ApiUrl)
}

func NewWithOptions(url string) *ApiClient {
    return &ApiClient{url: url}
}

func (c *ApiClient) AvailableTickers() ([]string, error) {
    data, err := c.apiCall("ticker/global/")
    if err != nil { return nil, err }

    var ti map[string]string
    err = json.Unmarshal(data, &ti)
    if err != nil { return nil, err }

    tl := make([]string, len(ti) - 1)
    i := 0
    for k := range ti {
        if k == "all" { continue }
        tl[i] = k
        i++
    }

    return tl, nil
}

func (c *ApiClient) Ignored() (map[string]string, error) {
    data, err := c.apiCall("ignored")
    if err != nil { return nil, err }

    var im map[string]string
    err = json.Unmarshal(data, &im)
    if err != nil { return nil, err }

    return im, nil
}

func (c *ApiClient) apiCall(endpoint string) ([]byte, error) {
    // Build URL
    url := fmt.Sprintf("%v/%v", c.url, endpoint)

    // Make request
    resp, err := http.Get(url)
    if err != nil { return nil, err }

    // Retrieve raw JSON response
    var body []byte
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil { return nil, err }
    defer resp.Body.Close()

    // Process API-level error conditions
    if resp.StatusCode != 200 {
        return nil, errors.New(string(body))
    }

    return body, nil
}
