package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Response is the main structure for the JSON response
type Response struct {
	Code    int        `json:"code"`
	Data    TickerData `json:"data"`
	Message string     `json:"message"`
}

// Data holds the inner data structure
type TickerData struct {
	Date   time.Time `json:"date"`
	Ticker Ticker    `json:"ticker"`
}

// Ticker holds the ticker data from JSON
type Ticker struct {
	Open int     `json:"open"`
	Last float64 `json:"last"`
	Past int64   `json:"past"`
	Vol  int     `json:"vol"`
}

// Custom UnmarshalJSON for the Ticker struct
func (t *Ticker) UnmarshalJSON(b []byte) error {
	type Alias Ticker
	aux := &struct {
		Last string `json:"last"` // Last is originally a string
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}

	// Convert the Last field from string to float64
	lastFloat, err := strconv.ParseFloat(aux.Last, 64)
	if err != nil {
		return err
	}
	t.Last = lastFloat

	return nil
}

// Custom UnmarshalJSON for the Data struct
func (d *TickerData) UnmarshalJSON(b []byte) error {
	type Alias TickerData
	aux := &struct {
		Date int64 `json:"date"` // Date is originally an int64 (timestamp)
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}

	// Convert timestamp to time.Time
	d.Date = time.Unix(0, aux.Date*int64(time.Millisecond))
	return nil
}

func main() {
	// Example JSON input
	jsonInput := `{
		"code": 0,
		"data": {
			"date": 1513865441609,
			"ticker": {
				"open": 10,
				"last": "10.15456456",
				"past": 21315456456,
				"vol": 110
			}
		},
		"message" : "Ok"
	}`

	var response Response
	if err := json.Unmarshal([]byte(jsonInput), &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println("Code:", response.Code)
	fmt.Println("Date:", response.Data.Date)
	fmt.Println("Ticker Open:", response.Data.Ticker.Open)
	fmt.Println("Ticker Last:", response.Data.Ticker.Last)
	fmt.Println("Ticker Past:", response.Data.Ticker.Past)
	fmt.Println("Ticker Vol:", response.Data.Ticker.Vol)
	fmt.Println("Message:", response.Message)
}
