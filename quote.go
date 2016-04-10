package quote

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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

// Price - stucture for historical price data
type Price struct {
	Symbol string      `json:"symbol"`
	Date   []time.Time `json:"date"`
	Open   []float64   `json:"open"`
	High   []float64   `json:"high"`
	Low    []float64   `json:"low"`
	Close  []float64   `json:"close"`
	Volume []float64   `json:"volume"`
}

// Prices - an array of Price
type Prices []Price

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseDTString(dt string) time.Time {
	t, err := time.Parse("2006-01-02 15:04", dt+"0000-01-01 00:00"[len(dt):])
	check(err)
	return t
}

// CSV - convert Price structure to csv string
func (p *Price) CSV() string {

	var buffer bytes.Buffer

	buffer.WriteString("datetime,open,high,low,close,volume\n")

	for bar := range p.Close {
		str := fmt.Sprintf("%s,%.2f,%.2f,%.2f,%.2f,%.0f\n",
			p.Date[bar].Format("2006-01-02 15:04"), p.Open[bar], p.High[bar], p.Low[bar], p.Close[bar], p.Volume[bar])
		buffer.WriteString(str)
	}

	return buffer.String()
}

// WriteCSV - write Price struct to csv file
func (p *Price) WriteCSV(filename string) {
	csv := p.CSV()
	ba := []byte(csv)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// NewPriceCSV - parse csv quote string into Price structure
func NewPriceCSV(stringOrFilename string) Price {

	var csv string
	if _, err := os.Stat(stringOrFilename); err == nil {
		raw, err := ioutil.ReadFile(stringOrFilename)
		check(err)
		csv = string(raw)
	} else {
		csv = stringOrFilename
	}

	p := Price{}
	tmp := strings.Split(csv, "\n")
	numrows := len(tmp) - 1
	p.Date = make([]time.Time, numrows-1)
	p.Open = make([]float64, numrows-1)
	p.High = make([]float64, numrows-1)
	p.Low = make([]float64, numrows-1)
	p.Close = make([]float64, numrows-1)
	p.Volume = make([]float64, numrows-1)

	for row, bar := 1, 0; row < numrows; row, bar = row+1, bar+1 {
		line := strings.Split(tmp[row], ",")
		p.Date[bar], _ = time.Parse("2006-01-02 15:04", line[0])
		p.Open[bar], _ = strconv.ParseFloat(line[1], 64)
		p.High[bar], _ = strconv.ParseFloat(line[2], 64)
		p.Low[bar], _ = strconv.ParseFloat(line[3], 64)
		p.Close[bar], _ = strconv.ParseFloat(line[4], 64)
		p.Volume[bar], _ = strconv.ParseFloat(line[5], 64)
	}
	return p
}

// JSON - convert Price struct to json string
func (p Price) JSON(indent bool) string {
	var j []byte
	if indent {
		j, _ = json.MarshalIndent(p, "", "  ")
	} else {
		j, _ = json.Marshal(p)
	}
	return string(j)
}

// WriteJSON - write Price struct to json file
func (p Price) WriteJSON(filename string, indent bool) {
	json := p.JSON(indent)
	ba := []byte(json)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// NewPriceJSON - parse json quote string into Price structure
func NewPriceJSON(stringOrFilename string) Price {

	var jsn []byte
	if _, err := os.Stat(stringOrFilename); err == nil {
		raw, err := ioutil.ReadFile(stringOrFilename)
		check(err)
		jsn = raw
	} else {
		jsn = []byte(stringOrFilename)
	}

	price := Price{}
	err := json.Unmarshal(jsn, &price)
	check(err)

	return price
}

// CSV - convert Prices structure to csv string
func (p Prices) CSV() string {

	var buffer bytes.Buffer

	buffer.WriteString("symbol,datetime,open,high,low,close,volume\n")

	for sym := 0; sym < len(p); sym++ {
		price := p[sym]
		for bar := range price.Close {
			str := fmt.Sprintf("%s,%s,%.2f,%.2f,%.2f,%.2f,%.0f\n",
				price.Symbol, price.Date[bar].Format("2006-01-02 15:04"), price.Open[bar], price.High[bar], price.Low[bar], price.Close[bar], price.Volume[bar])
			buffer.WriteString(str)
		}
	}

	return buffer.String()
}

// WriteCSV - write Prices structure to file
func (p Prices) WriteCSV(filename string) {
	csv := p.CSV()
	ba := []byte(csv)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// NewPricesCSV - parse csv quote string into Prices array
func NewPricesCSV(stringOrFilename string) Prices {

	var csv string
	if _, err := os.Stat(stringOrFilename); err == nil {
		raw, err := ioutil.ReadFile(stringOrFilename)
		check(err)
		csv = string(raw)
	} else {
		csv = stringOrFilename
	}

	prices := Prices{}

	tmp := strings.Split(csv, "\n")
	numrows := len(tmp) - 1

	var index = make(map[string]int)
	for idx := 1; idx < numrows; idx++ {
		sym := strings.Split(tmp[idx], ",")[0]
		index[sym]++
	}

	row := 1
	for sym, len := range index {
		p := Price{}
		p.Symbol = sym
		p.Date = make([]time.Time, len)
		p.Open = make([]float64, len)
		p.High = make([]float64, len)
		p.Low = make([]float64, len)
		p.Close = make([]float64, len)
		p.Volume = make([]float64, len)
		for bar := 0; bar < len; bar++ {
			line := strings.Split(tmp[row], ",")
			p.Date[bar], _ = time.Parse("2006-01-02 15:04", line[1])
			p.Open[bar], _ = strconv.ParseFloat(line[2], 64)
			p.High[bar], _ = strconv.ParseFloat(line[3], 64)
			p.Low[bar], _ = strconv.ParseFloat(line[4], 64)
			p.Close[bar], _ = strconv.ParseFloat(line[5], 64)
			p.Volume[bar], _ = strconv.ParseFloat(line[6], 64)
			row++
		}
		prices = append(prices, p)
	}
	return prices
}

// JSON - convert Prices struct to json string
func (p Prices) JSON(indent bool) string {
	var j []byte
	if indent {
		j, _ = json.MarshalIndent(p, "", "  ")
	} else {
		j, _ = json.Marshal(p)
	}
	return string(j)
}

// WriteJSON - write Price struct to json file
func (p Prices) WriteJSON(filename string, indent bool) {
	json := p.JSON(indent)
	ba := []byte(json)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// NewPricesJSON - parse json quote string into Price structure
func NewPricesJSON(stringOrFilename string) Prices {

	var jsn []byte
	if _, err := os.Stat(stringOrFilename); err == nil {
		raw, err := ioutil.ReadFile(stringOrFilename)
		check(err)
		jsn = raw
	} else {
		jsn = []byte(stringOrFilename)
	}

	prices := Prices{}
	err := json.Unmarshal(jsn, &prices)
	check(err)

	return prices
}

// NewYahoo - Yahoo historical prices for a symbol
func NewYahoo(symbol, startDate, endDate string, period Period, adjustPrice bool) (Price, error) {

	from := parseDTString(startDate)

	var to time.Time
	if endDate == "" {
		to = time.Now()
	} else {
		to = parseDTString(endDate)
	}

	price := Price{Symbol: symbol}

	url := fmt.Sprintf(
		"http://ichart.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%d&e=%d&f=%d&g=%s&ignore=.csv",
		symbol,
		from.Month()-1, from.Day(), from.Year(),
		to.Month()-1, to.Day(), to.Year(),
		period)

	resp, err := http.Get(url)
	if err != nil {
		return price, err
	}
	defer resp.Body.Close()

	var csvdata [][]string
	reader := csv.NewReader(resp.Body)
	csvdata, err = reader.ReadAll()
	if err != nil {
		return price, err
	}

	numrows := len(csvdata) - 1
	price.Date = make([]time.Time, numrows)
	price.Open = make([]float64, numrows)
	price.High = make([]float64, numrows)
	price.Low = make([]float64, numrows)
	price.Close = make([]float64, numrows)
	price.Volume = make([]float64, numrows)

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
		if adjustPrice {
			factor = a / c
		}

		// Append to prices
		bar := numrows - row // reverse the order
		price.Date[bar] = d
		price.Open[bar] = o * factor
		price.High[bar] = h * factor
		price.Low[bar] = l * factor
		price.Close[bar] = c * factor
		price.Volume[bar] = v

	}

	return price, nil
}

// NewYahooPrices - create a list of prices from symbols in file
func NewYahooPrices(filename, startDate, endDate string, period Period, adjustPrice bool) (Prices, error) {

	prices := Prices{}
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		sym := scanner.Text()
		price, _ := NewYahoo(sym, startDate, endDate, period, adjustPrice)
		prices = append(prices, price)
	}
	return prices, nil
}
