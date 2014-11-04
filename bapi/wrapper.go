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

type TickerData struct {
    // Average24h is not available when fetching all tickers in bulk through
    // Tickers()
    Average24h      json.Number     `json:"24h_avg"`
    Ask             json.Number     `json:"ask"`
    Bid             json.Number     `json:"bid"`
    Last            json.Number     `json:"last"`
    Timestamp       string          `json:"timestamp"`
    VolumeBTC       json.Number     `json:"volume_btc"`
    VolumePercent   json.Number     `json:"volume_percent"`
}

type AllTickerData struct {
    Tickers         map[string]TickerData
    Timestamp       string
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

func (c *ApiClient) Ticker(symbol string) (*TickerData, error) {
    data, err := c.apiCall("ticker/global/" + symbol)
    if err != nil { return nil, err }

    var td TickerData
    err = json.Unmarshal(data, &td)
    if err != nil { return nil, err }

    return &td, nil
}

func (c *ApiClient) Tickers() (*AllTickerData, error) {
    data, err := c.apiCall("ticker/global/all")
    if err != nil { return nil, err }

    // The API returns a nice map of symbols to TickerData, plus a timestamp...
    var td map[string]json.RawMessage
    err = json.Unmarshal(data, &td)
    if err != nil { return nil, err }

    var atd AllTickerData
    atd.Tickers = make(map[string]TickerData)
    for k, v := range td {
        if k == "timestamp" {
            err = json.Unmarshal(v, &atd.Timestamp)
            if err != nil { return nil, err }
            continue
        }

        var t TickerData
        err = json.Unmarshal(v, &t)
        if err != nil { return nil, err }
        atd.Tickers[k] = t
    }

    return &atd, nil
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
