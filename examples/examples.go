package main

import (
	"encoding/json"
	"fmt"
	"log"

	coinmarketcap "github.com/joemocquant/cmc-api"
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

func prettyPrintJson(msg interface{}) {

	jsonstr, err := json.MarshalIndent(msg, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(jsonstr))
}

func printTickers() {

	res, err := client.GetTickers()

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}

func printTickersLimit() {

	res, err := client.GetTickersLimit(10)

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}

func printTicker() {

	res, err := client.GetTicker("bitcoin")

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}

func printGlobalData() {

	res, err := client.GetGlobalData()

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}
