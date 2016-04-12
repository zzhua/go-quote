/*
Package quote is free quote downloader library and cli

Downloads daily/weekly/monthly historical price quotes from Yahoo
and daily/intraday data from Google

Copyright 2016 Mark Chenoweth
Licensed under terms of MIT license (see LICENSE)
*/
package quote

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Quote - stucture for historical price data
type Quote struct {
	Symbol string      `json:"symbol"`
	Date   []time.Time `json:"date"`
	Open   []float64   `json:"open"`
	High   []float64   `json:"high"`
	Low    []float64   `json:"low"`
	Close  []float64   `json:"close"`
	Volume []float64   `json:"volume"`
}

// Quotes - an array of historical price data
type Quotes []Quote

// Period - for quote history
type Period string

const (
	// Min1 - 1 Minute time period
	Min1 Period = "60"
	// Min5 - 5 Minute time period
	Min5 Period = "300"
	// Min15 - 15 Minute time period
	Min15 Period = "900"
	// Min30 - 30 Minute time period
	Min30 Period = "1800"
	// Min60 - 60 Minute time period
	Min60 Period = "3600"
	// Daily time period
	Daily Period = "d"
	// Weekly time period
	Weekly Period = "w"
	// Monthly time period
	Monthly Period = "m"
)

// NewQuote - new empty Quote struct
func NewQuote(symbol string, bars int) Quote {
	return Quote{
		Symbol: symbol,
		Date:   make([]time.Time, bars),
		Open:   make([]float64, bars),
		High:   make([]float64, bars),
		Low:    make([]float64, bars),
		Close:  make([]float64, bars),
		Volume: make([]float64, bars),
	}
}

// ParseDateString - parse a potentially partial date string to Time
func ParseDateString(dt string) time.Time {
	if dt == "" {
		return time.Now()
	}
	t, _ := time.Parse("2006-01-02 15:04", dt+"0000-01-01 00:00"[len(dt):])
	return t
}

// CSV - convert Quote structure to csv string
func (q Quote) CSV() string {
	var buffer bytes.Buffer
	buffer.WriteString("datetime,open,high,low,close,volume\n")
	for bar := range q.Close {
		str := fmt.Sprintf("%s,%.2f,%.2f,%.2f,%.2f,%.0f\n",
			q.Date[bar].Format("2006-01-02 15:04"), q.Open[bar], q.High[bar], q.Low[bar], q.Close[bar], q.Volume[bar])
		buffer.WriteString(str)
	}
	return buffer.String()
}

// WriteCSV - write Quote struct to csv file
func (q Quote) WriteCSV(filename string) error {
	if filename == "" {
		filename = q.Symbol + ".csv"
	}
	csv := q.CSV()
	return ioutil.WriteFile(filename, []byte(csv), 0644)
}

// NewQuoteFromCSV - parse csv quote string into Quote structure
func NewQuoteFromCSV(symbol, csv string) (Quote, error) {

	tmp := strings.Split(csv, "\n")
	numrows := len(tmp) - 1
	q := NewQuote("", numrows-1)

	for row, bar := 1, 0; row < numrows; row, bar = row+1, bar+1 {
		line := strings.Split(tmp[row], ",")
		q.Date[bar], _ = time.Parse("2006-01-02 15:04", line[0])
		q.Open[bar], _ = strconv.ParseFloat(line[1], 64)
		q.High[bar], _ = strconv.ParseFloat(line[2], 64)
		q.Low[bar], _ = strconv.ParseFloat(line[3], 64)
		q.Close[bar], _ = strconv.ParseFloat(line[4], 64)
		q.Volume[bar], _ = strconv.ParseFloat(line[5], 64)
	}
	return q, nil
}

// NewQuoteFromCSVFile - parse csv quote file into Quote structure
func NewQuoteFromCSVFile(symbol, filename string) (Quote, error) {
	csv, err := ioutil.ReadFile(filename)
	if err != nil {
		return NewQuote("", 0), err
	}
	return NewQuoteFromCSV(symbol, string(csv))
}

// JSON - convert Quote struct to json string
func (q Quote) JSON(indent bool) string {
	var j []byte
	if indent {
		j, _ = json.MarshalIndent(q, "", "  ")
	} else {
		j, _ = json.Marshal(q)
	}
	return string(j)
}

// WriteJSON - write Quote struct to json file
func (q Quote) WriteJSON(filename string, indent bool) error {
	if filename == "" {
		filename = q.Symbol + ".json"
	}
	json := q.JSON(indent)
	return ioutil.WriteFile(filename, []byte(json), 0644)

}

// NewQuoteFromJSON - parse json quote string into Quote structure
func NewQuoteFromJSON(jsn string) (Quote, error) {
	q := Quote{}
	err := json.Unmarshal([]byte(jsn), &q)
	if err != nil {
		return q, err
	}
	return q, nil
}

// NewQuoteFromJSONFile - parse json quote string into Quote structure
func NewQuoteFromJSONFile(filename string) (Quote, error) {
	jsn, err := ioutil.ReadFile(filename)
	if err != nil {
		return NewQuote("", 0), err
	}
	return NewQuoteFromJSON(string(jsn))
}

// CSV - convert Quotes structure to csv string
func (q Quotes) CSV() string {

	var buffer bytes.Buffer

	buffer.WriteString("symbol,datetime,open,high,low,close,volume\n")

	for sym := 0; sym < len(q); sym++ {
		quote := q[sym]
		for bar := range quote.Close {
			str := fmt.Sprintf("%s,%s,%.2f,%.2f,%.2f,%.2f,%.0f\n",
				quote.Symbol, quote.Date[bar].Format("2006-01-02 15:04"), quote.Open[bar], quote.High[bar], quote.Low[bar], quote.Close[bar], quote.Volume[bar])
			buffer.WriteString(str)
		}
	}

	return buffer.String()
}

// WriteCSV - write Quotes structure to file
func (q Quotes) WriteCSV(filename string) error {
	if filename == "" {
		filename = "quotes.csv"
	}
	csv := q.CSV()
	ba := []byte(csv)
	return ioutil.WriteFile(filename, ba, 0644)
}

// NewQuotesFromCSV - parse csv quote string into Quotes array
func NewQuotesFromCSV(csv string) (Quotes, error) {

	quotes := Quotes{}
	tmp := strings.Split(csv, "\n")
	numrows := len(tmp) - 1

	var index = make(map[string]int)
	for idx := 1; idx < numrows; idx++ {
		sym := strings.Split(tmp[idx], ",")[0]
		index[sym]++
	}

	row := 1
	for sym, len := range index {
		q := NewQuote(sym, len)
		for bar := 0; bar < len; bar++ {
			line := strings.Split(tmp[row], ",")
			q.Date[bar], _ = time.Parse("2006-01-02 15:04", line[1])
			q.Open[bar], _ = strconv.ParseFloat(line[2], 64)
			q.High[bar], _ = strconv.ParseFloat(line[3], 64)
			q.Low[bar], _ = strconv.ParseFloat(line[4], 64)
			q.Close[bar], _ = strconv.ParseFloat(line[5], 64)
			q.Volume[bar], _ = strconv.ParseFloat(line[6], 64)
			row++
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}

// NewQuotesFromCSVFile - parse csv quote file into Quotes array
func NewQuotesFromCSVFile(filename string) (Quotes, error) {
	csv, err := ioutil.ReadFile(filename)
	if err != nil {
		return Quotes{}, err
	}
	return NewQuotesFromCSV(string(csv))
}

// JSON - convert Quotes struct to json string
func (q Quotes) JSON(indent bool) string {
	var j []byte
	if indent {
		j, _ = json.MarshalIndent(q, "", "  ")
	} else {
		j, _ = json.Marshal(q)
	}
	return string(j)
}

// WriteJSON - write Quote struct to json file
func (q Quotes) WriteJSON(filename string, indent bool) error {
	if filename == "" {
		filename = "quotes.json"
	}
	jsn := q.JSON(indent)
	return ioutil.WriteFile(filename, []byte(jsn), 0644)
}

// NewQuotesFromJSON - parse json quote string into Quote structure
func NewQuotesFromJSON(jsn string) (Quotes, error) {
	quotes := Quotes{}
	err := json.Unmarshal([]byte(jsn), &quotes)
	if err != nil {
		return quotes, err
	}
	return quotes, nil
}

// NewQuotesFromJSONFile - parse json quote string into Quote structure
func NewQuotesFromJSONFile(filename string) (Quotes, error) {
	jsn, err := ioutil.ReadFile(filename)
	if err != nil {
		return Quotes{}, err
	}
	return NewQuotesFromJSON(string(jsn))
}

// NewQuoteFromYahoo - Yahoo historical prices for a symbol
func NewQuoteFromYahoo(symbol, startDate, endDate string, period Period, adjustQuote bool) (Quote, error) {

	from := ParseDateString(startDate)
	to := ParseDateString(endDate)

	url := fmt.Sprintf(
		"http://ichart.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%d&e=%d&f=%d&g=%s&ignore=.csv",
		symbol,
		from.Month()-1, from.Day(), from.Year(),
		to.Month()-1, to.Day(), to.Year(),
		period)

	resp, err := http.Get(url)
	if err != nil {
		return NewQuote("", 0), err
	}
	defer resp.Body.Close()

	var csvdata [][]string
	reader := csv.NewReader(resp.Body)
	csvdata, err = reader.ReadAll()
	if err != nil {
		return NewQuote("", 0), err
	}

	numrows := len(csvdata) - 1
	quote := NewQuote(symbol, numrows)

	for row := 1; row < len(csvdata); row++ {

		// Parse row of data
		d, _ := time.Parse("2006-01-02", csvdata[row][0])
		o, _ := strconv.ParseFloat(csvdata[row][1], 64)
		h, _ := strconv.ParseFloat(csvdata[row][2], 64)
		l, _ := strconv.ParseFloat(csvdata[row][3], 64)
		c, _ := strconv.ParseFloat(csvdata[row][4], 64)
		v, _ := strconv.ParseFloat(csvdata[row][5], 64)
		a, _ := strconv.ParseFloat(csvdata[row][6], 64)

		// Adjustment factor
		factor := 1.0
		if adjustQuote {
			factor = a / c
		}

		// Append to quote
		bar := numrows - row // reverse the order
		quote.Date[bar] = d
		quote.Open[bar] = o * factor
		quote.High[bar] = h * factor
		quote.Low[bar] = l * factor
		quote.Close[bar] = c * factor
		quote.Volume[bar] = v

	}

	return quote, nil
}

// NewQuotesFromYahoo - create a list of prices from symbols in file
func NewQuotesFromYahoo(filename, startDate, endDate string, period Period, adjustQuote bool) (Quotes, error) {

	quotes := Quotes{}
	inFile, err := os.Open(filename)
	if err != nil {
		return quotes, err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		sym := scanner.Text()
		quote, err := NewQuoteFromYahoo(sym, startDate, endDate, period, adjustQuote)
		if err == nil {
			quotes = append(quotes, quote)
		} //TODO else log error
	}
	return quotes, nil
}

// NewQuotesFromYahooSyms - create a list of prices from symbols in string array
func NewQuotesFromYahooSyms(symbols []string, startDate, endDate string, period Period, adjustQuote bool) (Quotes, error) {

	quotes := Quotes{}
	for _, symbol := range symbols {
		quote, err := NewQuoteFromYahoo(symbol, startDate, endDate, period, adjustQuote)
		if err == nil {
			quotes = append(quotes, quote)
		} //TODO else log error
	}
	return quotes, nil
}

func googleDaily(symbol string, from, to time.Time) (Quote, error) {

	args := fmt.Sprintf(
		"http://www.google.com/finance/historical?q=%s&startdate=%s&enddate=%s&output=csv",
		symbol,
		url.QueryEscape(from.Format("Jan 2, 2006")),
		url.QueryEscape(to.Format("Jan 2, 2006")))

	resp, err := http.Get(args)
	if err != nil {
		return NewQuote("", 0), err
	}
	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	tmp := strings.Join(strings.Split(string(contents), "\n")[1:], "\n")
	reader := csv.NewReader(strings.NewReader(tmp))
	csvdata, err := reader.ReadAll()
	if err != nil {
		return NewQuote("", 0), err
	}

	numrows := len(csvdata)
	quote := NewQuote(symbol, numrows)

	for row := 0; row < numrows; row++ {
		bar := numrows - 1 - row // reverse the order
		quote.Date[bar], _ = time.Parse("2-Jan-06", csvdata[row][0])
		quote.Open[bar], _ = strconv.ParseFloat(csvdata[row][1], 64)
		quote.High[bar], _ = strconv.ParseFloat(csvdata[row][2], 64)
		quote.Low[bar], _ = strconv.ParseFloat(csvdata[row][3], 64)
		quote.Close[bar], _ = strconv.ParseFloat(csvdata[row][4], 64)
		quote.Volume[bar], _ = strconv.ParseFloat(csvdata[row][5], 64)
	}

	return quote, nil
}

func googleIntra(symbol string, from, to time.Time, period Period) (Quote, error) {

	args := fmt.Sprintf(
		"http://www.google.com/finance/getprices?q=%s&i=%s&p=10d&f=d,o,h,l,c,v",
		strings.ToUpper(symbol),
		period)

	resp, err := http.Get(args)
	if err != nil {
		return NewQuote("", 0), err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	tmp := strings.Join(strings.Split(string(contents), "\n")[7:], "\n")
	reader := csv.NewReader(strings.NewReader(tmp))
	csvdata, err := reader.ReadAll()
	if err != nil {
		return NewQuote("", 0), err
	}

	numrows := len(csvdata)
	quote := NewQuote(symbol, numrows)

	var day int64
	for row := 0; row < numrows; row++ {

		var offset int64
		z := csvdata[row][0]
		if z[0] == 'a' {
			day, _ = strconv.ParseInt(z[1:], 10, 64)
		} else {
			offset, _ = strconv.ParseInt(z, 10, 64)
		}

		seconds, _ := strconv.ParseInt(string(period), 10, 64)
		quote.Date[row] = time.Unix(day+(seconds*offset), 0)
		quote.Open[row], _ = strconv.ParseFloat(csvdata[row][4], 64)
		quote.High[row], _ = strconv.ParseFloat(csvdata[row][2], 64)
		quote.Low[row], _ = strconv.ParseFloat(csvdata[row][3], 64)
		quote.Close[row], _ = strconv.ParseFloat(csvdata[row][1], 64)
		quote.Volume[row], _ = strconv.ParseFloat(csvdata[row][5], 64)
	}
	return quote, nil
}

// NewQuoteFromGoogle - Google daily/intraday historical prices for a symbol
func NewQuoteFromGoogle(symbol, startDate, endDate string, period Period) (Quote, error) {

	from := ParseDateString(startDate)
	to := ParseDateString(endDate)

	if period == Daily {
		return googleDaily(symbol, from, to)
	}
	return googleIntra(symbol, from, to, period)
}

// NewQuotesFromGoogle - create a list of prices from symbols in file
func NewQuotesFromGoogle(filename, startDate, endDate string, period Period) (Quotes, error) {

	quotes := Quotes{}
	inFile, err := os.Open(filename)
	if err != nil {
		return quotes, err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		sym := scanner.Text()
		quote, err := NewQuoteFromGoogle(sym, startDate, endDate, period)
		if err == nil {
			quotes = append(quotes, quote)
		} //TODO else log error
	}
	return quotes, nil
}

// NewQuotesFromGoogleSyms - create a list of prices from symbols in string array
func NewQuotesFromGoogleSyms(symbols []string, startDate, endDate string, period Period) (Quotes, error) {

	quotes := Quotes{}
	for _, symbol := range symbols {
		quote, err := NewQuoteFromGoogle(symbol, startDate, endDate, period)
		if err == nil {
			quotes = append(quotes, quote)
		} //TODO else log error
	}
	return quotes, nil
}

// NewEtfList - download a list of etf symbols to an array of strings
func NewEtfList() ([]string, error) {

	var symbols []string

	// http://www.nasdaqtrader.com/trader.aspx?id=symboldirdefs
	c, _ := ftp.DialTimeout("ftp.nasdaqtrader.com:21", 5*time.Second)
	_ = c.Login("anonymous", "anonymous")
	_ = c.ChangeDir("symboldirectory")
	r, _ := c.Retr("otherlisted.txt")
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return symbols, err
	}
	defer r.Close()

	for _, line := range strings.Split(string(buf), "\n") {
		// ACT Symbol|Security Name|Exchange|CQS Symbol|ETF|Round Lot Size|Test Issue|NASDAQ Symbol
		cols := strings.Split(line, "|")
		if len(cols) > 5 && cols[4] == "Y" && cols[6] == "N" {
			symbols = append(symbols, strings.ToLower(cols[0]))
		}
	}
	sort.Strings(symbols)
	return symbols, nil
}

// NewEtfFile - download a list of etf symbols to a file
func NewEtfFile(filename string) error {
	if filename == "" {
		filename = "etf.txt"
	}
	etfs, err := NewEtfList()
	if err != nil {
		return err
	}
	ba := []byte(strings.Join(etfs, "\n"))
	return ioutil.WriteFile(filename, ba, 0644)
}

func validExchange(exchange string) bool {
	return exchange == "nasdaq" || exchange == "nyse" || exchange == "amex"
}

// NewExchangeList - download a list of exchange symbols to an array of strings
func NewExchangeList(exchange string) ([]string, error) {

	var symbols []string
	if !validExchange(exchange) {
		return symbols, fmt.Errorf("invalid exchange")
	}

	url := fmt.Sprintf(
		"http://www.nasdaq.com/screening/companies-by-name.aspx?letter=0&exchange=%s&render=download",
		exchange)

	resp, err := http.Get(url)
	if err != nil {
		return symbols, err
	}
	defer resp.Body.Close()

	var csvdata [][]string
	reader := csv.NewReader(resp.Body)
	csvdata, err = reader.ReadAll()
	if err != nil {
		return symbols, err
	}

	r, _ := regexp.Compile("^[a-z]+$")
	for row := 1; row < len(csvdata); row++ {
		sym := strings.TrimSpace(strings.ToLower(csvdata[row][0]))
		if r.MatchString(sym) {
			symbols = append(symbols, sym)
		}
	}
	sort.Strings(symbols)
	return symbols, nil
}

// NewExchangeFile - download a list of exchange symbols to a file
func NewExchangeFile(exchange, filename string) error {

	if !validExchange(exchange) {
		return fmt.Errorf("invalid exchange")
	}

	// default filename
	if filename == "" {
		filename = exchange + ".txt"
	}

	syms, err := NewExchangeList(exchange)
	if err != nil {
		return err
	}
	ba := []byte(strings.Join(syms, "\n"))
	return ioutil.WriteFile(filename, ba, 0644)
}
