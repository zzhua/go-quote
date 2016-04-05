package quote

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	// Historical time periods:
	Daily   = "d"
	Weekly  = "w"
	Monthly = "m"
	Yearly  = "y"
)

// Price type that is used for a single day
type price struct {
	Date     time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
	AdjClose float64
}

// Price type that is used for historical price data.
type PriceH struct {
	Date   []time.Time `json:"date,omitempty"`
	Open   []float64   `json:"open,omitempty"`
	High   []float64   `json:"high,omitempty"`
	Low    []float64   `json:"low,omitempty"`
	Close  []float64   `json:"close,omitempty"`
	Volume []float64   `json:"volume,omitempty"`
	Adj    []float64   `json:"adj,omitempty"`
}

func (p *PriceH) Bar(at int) string {
	return fmt.Sprintf("%s,%.2f,%.2f,%.2f,%.2f,%.0f,%.2f", p.Date[at].Format("2006-01-02"), p.Open[at], p.High[at], p.Low[at], p.Close[at], p.Volume[at], p.Adj[at])
}

// Get historical prices for the stock.
func GetDailyHistory(symbol string, from, to time.Time, adjustPrice bool) (PriceH, error) {
	var prices PriceH
	// Create URL with daily frequency of data.
	data, err := loadHistoricalPrice(symbol, from, to, Daily)
	if err != nil {
		return prices, err
	}

	prices, err = parseHistorical(data, adjustPrice)
	if err != nil {
		return prices, err
	}

	return prices, nil
}

// Get stock price history for number of years backwards.
func HistoryForYears(symbol string, years int, period string, adjustPrice bool) (PriceH, error) {
	var prices PriceH
	duration := time.Duration(int(time.Hour) * 24 * 365 * years)
	to := time.Now()
	from := to.Add(-duration)

	data, err := loadHistoricalPrice(symbol, from, to, period)
	if err != nil {
		return prices, err
	}

	prices, err = parseHistorical(data, adjustPrice)
	if err != nil {
		return prices, err
	}

	return prices, nil
}

// Load historical data price.
func loadHistoricalPrice(symbol string, from, to time.Time, period string) ([][]string, error) {
	url := stockHistoryURL(symbol, from, to, period)
	var data [][]string

	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	data, err = reader.ReadAll()
	if err != nil {
		return data, err
	}

	return data, nil
}

// Generate request URL for historicat stock data.
func stockHistoryURL(symbol string, from, to time.Time, frequency string) string {

	// Historical data URL with params:
	// s - symbol
	// a - from month (zero based)
	// b - from day
	// c - from year
	// d - to month (zero based)
	// e - to day
	// f - to year
	// g - period frequence (d - daily, w - weekly, m - monthly, y -yearly)
	yurl := "http://ichart.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%d&e=%d&f=%d&g=%s&ignore=.csv"

	// From date
	fMonth := (from.Month() - 1) // Need to subtract 1 because months in query is 0 based.
	fDay := from.Day()
	fYear := from.Year()
	// To date
	tMonth := (to.Month() - 1)
	tDay := to.Day()
	tYear := to.Year()

	url := fmt.Sprintf(
		yurl,
		symbol,
		fMonth,
		fDay,
		fYear,
		tMonth,
		tDay,
		tYear,
		frequency)

	return url
}

// Parse collection of historical prices.
func parseHistorical(data [][]string, adjustPrice bool) (PriceH, error) {
	// This is the list of prices with allocated space. Length of space should
	// subtracted by 1 because the first row of data is title.
	dl := len(data) - 1
	prices := PriceH{}
	prices.Date = make([]time.Time, dl)
	prices.Open = make([]float64, dl)
	prices.High = make([]float64, dl)
	prices.Low = make([]float64, dl)
	prices.Close = make([]float64, dl)
	prices.Volume = make([]float64, dl)
	prices.Adj = make([]float64, dl)

	// We need to leave the first row, because it contains title of columns.
	for k, v := range data {
		if k == 0 {
			continue
		}
		// Parse row of data into PriceH type and append it to collection of prices.
		p, err := parseHistoricalRow(v)
		if err != nil {
			return prices, err
		}

		// dl because we are reversing the returned data
		// (k - 1) because we remove header from the list so index should be reduced by one.
		prices.Date[dl-k] = p.Date
		if adjustPrice {
			factor := p.AdjClose / p.Close
			prices.Open[dl-k] = p.Open * factor
			prices.High[dl-k] = p.High * factor
			prices.Low[dl-k] = p.Low * factor
			prices.Close[dl-k] = p.Close * factor
			prices.Volume[dl-k] = p.Volume
			prices.Adj[dl-k] = factor
		} else {
			prices.Open[dl-k] = p.Open
			prices.High[dl-k] = p.High
			prices.Low[dl-k] = p.Low
			prices.Close[dl-k] = p.Close
			prices.Volume[dl-k] = p.Volume
			prices.Adj[dl-k] = p.AdjClose
		}
	}
	return prices, nil
}

// Parse data row that comes from historical data. Data row contains
// 7 columns:
// 0 - Date
// 1 - Open
// 2 - High
// 3 - Low
// 4 - Close
// 5 - Volume
// 6 - Adj Close
// This function will return PriceH type that wraps all these columns.
func parseHistoricalRow(data []string) (price, error) {
	p := price{}

	// Parse date.
	d, err := time.Parse("2006-01-02", data[0])
	if err != nil {
		return p, err
	}

	p.Date = d
	p.Open, _ = strconv.ParseFloat(data[1], 64)
	p.High, _ = strconv.ParseFloat(data[2], 64)
	p.Low, _ = strconv.ParseFloat(data[3], 64)
	p.Close, _ = strconv.ParseFloat(data[4], 64)
	p.Volume, _ = strconv.ParseFloat(data[5], 64)
	p.AdjClose, _ = strconv.ParseFloat(data[6], 64)

	return p, nil
}
