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

type Ticker struct {
    // Average24h is not available when fetching all tickers in bulk through
    // GlobalTickers()
    Average24h      json.Number     `json:"24h_avg"`
    Ask             json.Number     `json:"ask"`
    Bid             json.Number     `json:"bid"`
    Last            json.Number     `json:"last"`
    Timestamp       string          `json:"timestamp"`
    // Volume* only available for global tickers.
    VolumeBTC       json.Number     `json:"volume_btc"`
    VolumePercent   json.Number     `json:"volume_percent"`
    // TotalVolume is only available for market tickers.
    TotalVolume     json.Number     `json:"total_vol"`
}

type AllTickers struct {
    Tickers         map[string]Ticker
    Timestamp       string
}

type Exchange struct {
    DisplayURL      string          `json:"display_URL"`
    DisplayName     string          `json:"display_name"`
    Rates           ExchangeRates   `json:"rates"`
    Source          string          `json:"source"`
    VolumeBTC       json.Number     `json:"volume_btc"`
    VolumePercent   json.Number     `json:"volume_percent"`
}

type ExchangeRates struct {
    Ask             json.Number     `json:"ask"`
    Bid             json.Number     `json:"bid"`
    Last            json.Number     `json:"last"`
}

type ExchangeList struct {
    Exchanges       map[string]Exchange
    Timestamp       string
}

type AllExchanges struct {
    Exchanges       map[string]map[string]Exchange
    Timestamp       string
}

func New() *ApiClient {
    return NewWithOptions(ApiUrl)
}

func NewWithOptions(url string) *ApiClient {
    return &ApiClient{url: url}
}

func (c *ApiClient) GlobalTickerList() ([]string, error) {
    return c.index("ticker/global/")
}

func (c *ApiClient) MarketTickerList() ([]string, error) {
    return c.index("ticker/")
}

func (c *ApiClient) ExchangeList() ([]string, error) {
    return c.index("exchanges/")
}

func (c *ApiClient) index(endpoint string) ([]string, error) {
    data, err := c.apiCall(endpoint)
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

func (c *ApiClient) GlobalTicker(symbol string) (*Ticker, error) {
    return c.ticker("ticker/global/", symbol)
}

func (c *ApiClient) MarketTicker(symbol string) (*Ticker, error) {
    return c.ticker("ticker/", symbol)
}

func (c *ApiClient) ticker(endpoint string, symbol string) (*Ticker, error) {
    data, err := c.apiCall(endpoint + symbol)
    if err != nil { return nil, err }

    var t Ticker
    err = json.Unmarshal(data, &t)
    if err != nil { return nil, err }

    return &t, nil
}

func (c *ApiClient) GlobalTickers() (*AllTickers, error) {
    return c.tickers("ticker/global/all")
}

func (c *ApiClient) MarketTickers() (*AllTickers, error) {
    return c.tickers("ticker/all")
}

func (c *ApiClient) tickers(endpoint string) (*AllTickers, error) {
    data, err := c.apiCall(endpoint)
    if err != nil { return nil, err }

    // The API returns a nice map of symbols to Ticker, plus a timestamp...
    var td map[string]json.RawMessage
    err = json.Unmarshal(data, &td)
    if err != nil { return nil, err }

    var at AllTickers
    at.Tickers = make(map[string]Ticker)
    for k, v := range td {
        if k == "timestamp" {
            err = json.Unmarshal(v, &at.Timestamp)
            if err != nil { return nil, err }
            continue
        }

        var t Ticker
        err = json.Unmarshal(v, &t)
        if err != nil { return nil, err }
        at.Tickers[k] = t
    }

    return &at, nil
}

func (c *ApiClient) Exchanges(symbol string) (*ExchangeList, error) {
    data, err := c.apiCall("exchanges/" + symbol)
    if err != nil { return nil, err }

    // The API returns a nice map of names to Exchange, plus a timestamp...
    var ed map[string]json.RawMessage
    err = json.Unmarshal(data, &ed)
    if err != nil { return nil, err }

    var el ExchangeList
    el.Exchanges = make(map[string]Exchange)
    for k, v := range ed {
        if k == "timestamp" {
            err = json.Unmarshal(v, &el.Timestamp)
            if err != nil { return nil, err }
            continue
        }

        var e Exchange
        err = json.Unmarshal(v, &e)
        if err != nil { return nil, err }
        el.Exchanges[k] = e
    }

    return &el, nil
}

func (c *ApiClient) AllExchanges() (*AllExchanges, error) {
    data, err := c.apiCall("exchanges/all")
    if err != nil { return nil, err }

    // The API returns a nice map of symbols to Exchange, plus a timestamp...
    var ed map[string]json.RawMessage
    err = json.Unmarshal(data, &ed)
    if err != nil { return nil, err }

    var ae AllExchanges
    ae.Exchanges = make(map[string]map[string]Exchange)
    for k, v := range ed {
        if k == "timestamp" {
            err = json.Unmarshal(v, &ae.Timestamp)
            if err != nil { return nil, err }
            continue
        }

        var e map[string]Exchange
        err = json.Unmarshal(v, &e)
        if err != nil { return nil, err }
        ae.Exchanges[k] = e
    }

    return &ae, nil
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
