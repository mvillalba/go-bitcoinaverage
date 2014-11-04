package bapi

import (
    "encoding/json"
    "encoding/csv"
    "io/ioutil"
    "net/http"
    "strings"
    "errors"
    "bytes"
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

type HourlyHistoryRecord struct {
    DateTime        string
    High            json.Number
    Low             json.Number
    Average         json.Number
}

type DailyHistoryRecord struct {
    DateTime        string
    High            json.Number
    Low             json.Number
    Average         json.Number
    Volume          json.Number
}

type VolumeHistoryRecord struct {
    DateTime        string
    TotalVolume     json.Number
    Exchanges       map[string]ExchangeVolumeHistoryRecord
}

type ExchangeVolumeHistoryRecord struct {
    VolumeBTC       json.Number
    VolumePercent   json.Number
}

func New() *ApiClient {
    return NewWithOptions(ApiUrl)
}

func NewWithOptions(url string) *ApiClient {
    return &ApiClient{url: url}
}

func (c *ApiClient) GlobalTickerList() ([]string, error) {
    return c.index("ticker/global/", true)
}

func (c *ApiClient) MarketTickerList() ([]string, error) {
    return c.index("ticker/", true)
}

func (c *ApiClient) ExchangeList() ([]string, error) {
    return c.index("exchanges/", true)
}

func (c *ApiClient) HistoryList() ([]string, error) {
    return c.index("history/", false)
}

func (c *ApiClient) index(endpoint string, hasAll bool) ([]string, error) {
    data, err := c.apiCall(endpoint)
    if err != nil { return nil, err }

    var ti map[string]string
    err = json.Unmarshal(data, &ti)
    if err != nil { return nil, err }

    tllen := len(ti)
    if hasAll { tllen -= 1 }
    tl := make([]string, tllen)
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

func (c *ApiClient) HourlyHistory(symbol string) ([]HourlyHistoryRecord, error) {
    header, records, err := c.csvCall("history/" + symbol + "/per_hour_monthly_sliding_window.csv")
    if err != nil { return nil, err }

    rs := make([]HourlyHistoryRecord, len(records))
    for idx, record := range records {
        var r HourlyHistoryRecord

        for i, column := range header {
            switch column {
            case "datetime": r.DateTime = record[i];
            case "high": r.High = json.Number(record[i]);
            case "low": r.Low = json.Number(record[i]);
            case "average": r.Average = json.Number(record[i]);
            default: return nil, errors.New("got unexpected CSV columns.")
            }
        }

        rs[idx] = r
    }

    return rs, nil
}

func (c *ApiClient) DailyHistory(symbol string) ([]DailyHistoryRecord, error) {
    header, records, err := c.csvCall("history/" + symbol + "/per_day_all_time_history.csv")
    if err != nil { return nil, err }

    rs := make([]DailyHistoryRecord, len(records))
    for idx, record := range records {
        var r DailyHistoryRecord

        for i, column := range header {
            switch column {
            case "datetime": r.DateTime = record[i];
            case "high": r.High = json.Number(record[i]);
            case "low": r.Low = json.Number(record[i]);
            case "average": r.Average = json.Number(record[i]);
            case "volume": r.Volume = json.Number(record[i]);
            default: return nil, errors.New("got unexpected CSV columns.")
            }
        }

        rs[idx] = r
    }

    return rs, nil
}

func (c *ApiClient) VolumeHistory(symbol string) ([]VolumeHistoryRecord, error) {
    // Fetch CSV
    header, records, err := c.csvCall("history/" + symbol + "/volumes.csv")
    if err != nil { return nil, err }

    // Process as best we can
    var rs []VolumeHistoryRecord
    for _, record := range records {
        var r VolumeHistoryRecord
        r.Exchanges = make(map[string]ExchangeVolumeHistoryRecord)

        for i, column := range header {
            switch column {
            case "datetime": r.DateTime = record[i];
            case "total_vol": r.TotalVolume = json.Number(record[i]);
            default:
                m := strings.Split(column, " ")
                if len(m) != 2 { return nil, errors.New("got malformed CSV data.") }
                val := r.Exchanges[m[0]]
                if m[1] == "BTC" {
                    val.VolumeBTC = json.Number(record[i])
                } else {
                    val.VolumePercent = json.Number(record[i])
                }
                r.Exchanges[m[0]] = val
            }
        }

        // Filter bogus data (argggggggg......)
        for k, v := range r.Exchanges {
            if (v.VolumeBTC == "" || v.VolumePercent == "") ||
               (v.VolumeBTC == "0" && v.VolumePercent == "0") {
                   delete(r.Exchanges, k)
            }
        }

        if len(r.Exchanges) >= 1 {
            rs = append(rs, r)
        }
    }

    return rs, nil
}

func (c *ApiClient) csvCall(endpoint string) ([]string, [][]string, error) {
    // Fetch CSV
    data, err := c.apiCall(endpoint)
    if err != nil { return nil, nil, err }

    // Initialize CSV reader
    stream := bytes.NewReader(data)
    reader := csv.NewReader(stream)

    // Get CSV header
    var header []string
    header, err = reader.Read()
    if err != nil { return nil, nil, err }

    // Get CSV records
    var records [][]string
    records, err = reader.ReadAll()
    if err != nil { return nil, nil, err }

    return header, records, nil
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
