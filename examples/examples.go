package main

import (
	"log"
	"trading/coinmarketcap"
	"trading/poloniex"
)

var client *coinmarketcap.Client

// go run example.go
func main() {

	client = coinmarketcap.NewClient()

	// printTickers()
	// printTickersLimit()
	// printTicker()
	printGlobalData()
}

func printTickers() {

	res, err := client.GetTickers()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

func printTickersLimit() {

	res, err := client.GetTickersLimit(10)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

func printTicker() {

	res, err := client.GetTicker("bitcoin")

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

func printGlobalData() {

	res, err := client.GetGlobalData()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}
