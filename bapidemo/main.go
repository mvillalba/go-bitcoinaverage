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
    ignored(client)
    /*supportedCurrencies(client)
    currentPrice(client)
    currentPriceForCurrency(client)
    historical(client)
    historicalForYesterday(client)
    historicalForDates(client)*/
}

func version() {
    fmt.Println("BAPI Interface Version:", bapi.Version)
    fmt.Println("BAPI Interface Author:", bapi.Author)
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
