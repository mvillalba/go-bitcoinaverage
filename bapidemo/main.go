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
    availableTickers(client)
    ticker(client)
    tickers(client)
    ignored(client)
}

func version() {
    fmt.Println("BAPI Interface Version:", bapi.Version)
    fmt.Println("BAPI Interface Author:", bapi.Author)
}

func availableTickers(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List available tickers.")
    fmt.Println("=======================================")
    fmt.Println("")

    tl, err := client.AvailableTickers()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for _, t := range tl {
        fmt.Println(t)
    }
}

func ticker(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("Display GBP ticker data.")
    fmt.Println("=======================================")
    fmt.Println("")

    t, err := client.Ticker("GBP")
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

func tickers(client *bapi.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("Display all tickers.")
    fmt.Println("=======================================")
    fmt.Println("")

    tl, err := client.Tickers()
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
