package quote

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Prices - Stucture for historical price data
type Prices struct {
	Symbol string      `json:"symbol"`
	Date   []time.Time `json:"date"`
	Open   []float64   `json:"open"`
	High   []float64   `json:"high"`
	Low    []float64   `json:"low"`
	Close  []float64   `json:"close"`
	Volume []float64   `json:"volume"`
}

// Period - for quote history
type Period string

const (
	// Daily time period
	Daily Period = "d"
	// Weekly time period
	Weekly Period = "w"
	// Monthly time period
	Monthly Period = "m"
	// Yearly time period
	Yearly Period = "y"
)

// ToCSV - convert prices to csv
func (p *Prices) ToCSV(headerRow, symbolColumn bool) string {

	var buffer bytes.Buffer

	if headerRow {
		if symbolColumn {
			buffer.WriteString("symbol,")
		}
		buffer.WriteString("date,open,high,low,close,volume\n")
	}

	for bar := range p.Close {
		sc := ""
		if symbolColumn {
			sc = p.Symbol + ","
		}
		dt := fmt.Sprintf("%d-%02d-%02d", p.Date[bar].Year(), p.Date[bar].Month(), p.Date[bar].Day())
		str := fmt.Sprintf("%s%s,%.2f,%.2f,%.2f,%.2f,%.0f\n",
			sc, dt, p.Open[bar], p.High[bar], p.Low[bar], p.Close[bar], p.Volume[bar])
		buffer.WriteString(str)
	}

	return buffer.String()
}

// ToJSON - convert prices to json
func (p *Prices) ToJSON(indent bool) string {
	var j []byte
	if indent {
		j, _ = json.MarshalIndent(p, "", "  ")
	} else {
		j, _ = json.Marshal(p)
	}
	return string(j)
}

// NewYahoo - Yahoo historical prices for a symbol
func NewYahoo(symbol string, from, to time.Time, period Period, adjustPrice bool) (*Prices, error) {

	prices := &Prices{}
	prices.Symbol = symbol

	url := fmt.Sprintf(
		"http://ichart.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%d&e=%d&f=%d&g=%s&ignore=.csv",
		symbol,
		from.Month()-1, from.Day(), from.Year(),
		to.Month()-1, to.Day(), to.Year(),
		period)

	resp, err := http.Get(url)
	if err != nil {
		return prices, err
	}
	defer resp.Body.Close()

	var csvdata [][]string
	reader := csv.NewReader(resp.Body)
	csvdata, err = reader.ReadAll()
	if err != nil {
		return prices, err
	}

	numrows := len(csvdata) - 1
	prices.Date = make([]time.Time, numrows)
	prices.Open = make([]float64, numrows)
	prices.High = make([]float64, numrows)
	prices.Low = make([]float64, numrows)
	prices.Close = make([]float64, numrows)
	prices.Volume = make([]float64, numrows)

	for row := range csvdata {

		if row == 0 { // skip header
			continue
		}

		bar := numrows - row // reverse the order

		// Parse row of data and append it to collection of prices

		d, _ := time.Parse("2006-01-02", csvdata[row][0])
		o, _ := strconv.ParseFloat(csvdata[row][1], 64)
		h, _ := strconv.ParseFloat(csvdata[row][2], 64)
		l, _ := strconv.ParseFloat(csvdata[row][3], 64)
		c, _ := strconv.ParseFloat(csvdata[row][4], 64)
		v, _ := strconv.ParseFloat(csvdata[row][5], 64)
		a, _ := strconv.ParseFloat(csvdata[row][6], 64)

		prices.Date[bar] = d
		prices.Volume[bar] = v

		if adjustPrice {
			factor := a / c
			prices.Open[bar] = o * factor
			prices.High[bar] = h * factor
			prices.Low[bar] = l * factor
			prices.Close[bar] = c * factor
		} else {
			prices.Open[bar] = o
			prices.High[bar] = h
			prices.Low[bar] = l
			prices.Close[bar] = c
		}
	}

	return prices, nil
}

// NewYahooYears - get Yahoo stock price history for a number of years
func NewYahooYears(symbol string, years int, period Period, adjustPrice bool) (*Prices, error) {
	to := time.Now()
	from := to.Add(-time.Duration(int(time.Hour) * 24 * 365 * years))
	return NewYahoo(symbol, from, to, period, adjustPrice)
}
