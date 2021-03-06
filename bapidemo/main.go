package main

import (
    "github.com/mvillalba/go-bitcoinaverage/bapi"
    "fmt"
)

func main() {
    // Init ApiClient
    client := bapi.New()

    // Fun stuff
    version()
    globalTickerList(client)
    globalTicker(client)
    globalTickers(client)
    marketTickerList(client)
    marketTicker(client)
    marketTickers(client)
    exchangeList(client)
    exchanges(client)
    allExchanges(client)
    historyList(client)
    minutelyHistory(client)
    hourlyHistory(client)
    dailyHistory(client)
    volumeHistory(client)
    ignored(client)
}

func version() {
    fmt.Println("BAPI Interface Version:", bapi.Version)
    fmt.Println("BAPI Interface Author:", bapi.Author)
}

func globalTickerList(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List available global tickers.")
    fmt.Println("=======================================")
    fmt.Println("")

    tl, err := client.GlobalTickerList()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, t := range tl {
        fmt.Println(t)
    }
}

func globalTicker(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("Display GBP global ticker data.")
    fmt.Println("=======================================")
    fmt.Println("")

    t, err := client.GlobalTicker("GBP")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    fmt.Println("Average24h:", t.Average24h)
    fmt.Println("Ask:", t.Ask)
    fmt.Println("Bid:", t.Bid)
    fmt.Println("Last:", t.Last)
    fmt.Println("Timestamp:", t.Timestamp)
    fmt.Println("VolumeBTC:", t.VolumeBTC)
    fmt.Println("VolumePercent:", t.VolumePercent)
}

func globalTickers(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("Display all global tickers.")
    fmt.Println("=======================================")
    fmt.Println("")

    tl, err := client.GlobalTickers()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    fmt.Println("Timestamp:", tl.Timestamp)
    fmt.Println("")

    for s, t := range tl.Tickers {
        fmt.Println(s + ":")
        fmt.Println("  Ask:", t.Ask)
        fmt.Println("  Bid:", t.Bid)
        fmt.Println("  Last:", t.Last)
        fmt.Println("  Timestamp:", t.Timestamp)
        fmt.Println("  VolumeBTC:", t.VolumeBTC)
        fmt.Println("  VolumePercent:", t.VolumePercent)
    }
}

func marketTickerList(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List available market tickers.")
    fmt.Println("=======================================")
    fmt.Println("")

    tl, err := client.MarketTickerList()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, t := range tl {
        fmt.Println(t)
    }
}

func marketTicker(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("Display AUD market ticker data.")
    fmt.Println("=======================================")
    fmt.Println("")

    t, err := client.MarketTicker("AUD")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    fmt.Println("Average24h:", t.Average24h)
    fmt.Println("Ask:", t.Ask)
    fmt.Println("Bid:", t.Bid)
    fmt.Println("Last:", t.Last)
    fmt.Println("Timestamp:", t.Timestamp)
    fmt.Println("TotalVolume:", t.TotalVolume)
}

func marketTickers(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("Display all market tickers.")
    fmt.Println("=======================================")
    fmt.Println("")

    tl, err := client.MarketTickers()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    fmt.Println("Timestamp:", tl.Timestamp)
    fmt.Println("")

    for s, t := range tl.Tickers {
        fmt.Println(s + ":")
        fmt.Println("  Average24h:", t.Average24h)
        fmt.Println("  Ask:", t.Ask)
        fmt.Println("  Bid:", t.Bid)
        fmt.Println("  Last:", t.Last)
        fmt.Println("  Timestamp:", t.Timestamp)
        fmt.Println("  TotalVolume:", t.TotalVolume)
    }
}

func exchangeList(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List symbols for which there is")
    fmt.Println("exchange-specific data available.")
    fmt.Println("=======================================")
    fmt.Println("")

    tl, err := client.ExchangeList()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, t := range tl {
        fmt.Println(t)
    }
}

func exchanges(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List all exchanges for symbols USD.")
    fmt.Println("=======================================")
    fmt.Println("")

    el, err := client.Exchanges("USD")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    fmt.Println("Timestamp:", el.Timestamp)
    fmt.Println("")

    for k, v := range el.Exchanges {
        fmt.Println(k + ":")
        fmt.Println("  DisplayURL:", v.DisplayURL)
        fmt.Println("  DisplayName:", v.DisplayName)
        fmt.Println("  Rates / Ask:", v.Rates.Ask)
        fmt.Println("  Rates / Bid:", v.Rates.Bid)
        fmt.Println("  Rates / Last:", v.Rates.Last)
        fmt.Println("  Source:", v.Source)
        fmt.Println("  VolumeBTC:", v.VolumeBTC)
        fmt.Println("  VolumePercent:", v.VolumePercent)
    }
}

func allExchanges(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List all exchanges for all symbols.")
    fmt.Println("=======================================")
    fmt.Println("")

    el, err := client.AllExchanges()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    fmt.Println("Timestamp:", el.Timestamp)
    fmt.Println("")

    for k, v := range el.Exchanges {
        fmt.Println(k + ":")
        for kk, vv := range v {
            fmt.Println("  " + kk + ":")
            fmt.Println("    DisplayURL:", vv.DisplayURL)
            fmt.Println("    DisplayName:", vv.DisplayName)
            fmt.Println("    Rates / Ask", vv.Rates.Ask)
            fmt.Println("    Rates / Bid:", vv.Rates.Bid)
            fmt.Println("    Rates / Last:", vv.Rates.Last)
            fmt.Println("    Source:", vv.Source)
            fmt.Println("    VolumeBTC:", vv.VolumeBTC)
            fmt.Println("    VolumePercent:", vv.VolumePercent)
        }
    }
}

func historyList(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List symbols for which there is")
    fmt.Println("historical data available.")
    fmt.Println("=======================================")
    fmt.Println("")

    hl, err := client.HistoryList()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, s := range hl {
        fmt.Println(s)
    }
}

func minutelyHistory(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List minute by minute sliding window")
    fmt.Println("history data for the last 24 hours for")
    fmt.Println("symbol SGD.")
    fmt.Println("=======================================")

    vh, err := client.MinutelyHistory("SGD")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, r := range vh {
        fmt.Println("")
        fmt.Println("DateTime:", r.DateTime)
        fmt.Println("Average:", r.Average)
    }
}

func hourlyHistory(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List hourly sliding window history data")
    fmt.Println("for the last 30 days for symbol HKD.")
    fmt.Println("=======================================")

    vh, err := client.HourlyHistory("HKD")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, r := range vh {
        fmt.Println("")
        fmt.Println("DateTime:", r.DateTime)
        fmt.Println("High:", r.High)
        fmt.Println("Low:", r.Low)
        fmt.Println("Average:", r.Average)
    }
}

func dailyHistory(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List daily history data for symbol AUD.")
    fmt.Println("=======================================")

    vh, err := client.DailyHistory("AUD")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, r := range vh {
        fmt.Println("")
        fmt.Println("DateTime:", r.DateTime)
        fmt.Println("High:", r.High)
        fmt.Println("Low:", r.Low)
        fmt.Println("Average:", r.Average)
        fmt.Println("Volume:", r.Volume)
    }
}

func volumeHistory(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List volume history data for symbol")
    fmt.Println("NZD.")
    fmt.Println("=======================================")

    vh, err := client.VolumeHistory("NZD")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, r := range vh {
        fmt.Println("")
        fmt.Println("DateTime:", r.DateTime)
        fmt.Println("TotalVolume:", r.TotalVolume)
        for k, v := range r.Exchanges {
            fmt.Println(k + ":")
            fmt.Println("  VolumeBTC:", v.VolumeBTC)
            fmt.Println("  VolumePercent:", v.VolumePercent)
        }
    }
}

func ignored(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List all ignored exchanges.")
    fmt.Println("=======================================")
    fmt.Println("")

    im, err := client.Ignored()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for k := range im {
        fmt.Println(k + ":", im[k])
    }
}
