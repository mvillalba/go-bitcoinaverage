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
    availableTickers(client)
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
