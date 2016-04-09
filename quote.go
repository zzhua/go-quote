package quote

// p := quote.NewYahoo("spy","2000","",quote.Daily,true)
// p := quote.NewYahoo("spy","2000","2010",quote.Daily,true)
// p := quote.NewYahoo("spy","2000-02-01","2010-12-31",quote.Daily,true)
// p := quote.NewYahooYears("spy",5,quote.Daily,true)
// csv := p.ToCsv(true,true)
// p1 := quote.PricesFromCSV(csv)
// p.WriteCSV("spy.csv")
// p1 := quote.ReadPrices("spy.csv")
// s := quote.ReadSymbols("symbols.csv")
// s := quote.NewYahooSymbols("etf.txt")

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

// Symbols - an array of Prices
type Symbols []Prices

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

// ToCSV - convert Prices structure to csv string
func (p *Prices) ToCSV() string {

	var buffer bytes.Buffer

	buffer.WriteString("date,open,high,low,close,volume\n")

	for bar := range p.Close {
		str := fmt.Sprintf("%s,%.2f,%.2f,%.2f,%.2f,%.0f\n",
			p.Date[bar].Format("2006-01-02 15:04"), p.Open[bar], p.High[bar], p.Low[bar], p.Close[bar], p.Volume[bar])
		buffer.WriteString(str)
	}

	return buffer.String()
}

// WriteCSV - write Prices struct to csv file
func (p *Prices) WriteCSV(filename string) {
	csv := p.ToCSV()
	ba := []byte(csv)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// PricesFromCSV - parse csv quote string into Prices structure
func PricesFromCSV(csv string) Prices {

	p := Prices{}
	tmp := strings.Split(csv, "\n")
	numrows := len(tmp) - 1
	p.Date = make([]time.Time, numrows)
	p.Open = make([]float64, numrows)
	p.High = make([]float64, numrows)
	p.Low = make([]float64, numrows)
	p.Close = make([]float64, numrows)
	p.Volume = make([]float64, numrows)

	for row, bar := 1, 0; row < numrows; row, bar = row+1, bar+1 {
		line := strings.Split(tmp[row], ",")
		p.Date[bar], _ = time.Parse("2006-01-02", line[0])
		p.Open[bar], _ = strconv.ParseFloat(line[1], 64)
		p.High[bar], _ = strconv.ParseFloat(line[2], 64)
		p.Low[bar], _ = strconv.ParseFloat(line[3], 64)
		p.Close[bar], _ = strconv.ParseFloat(line[4], 64)
		p.Volume[bar], _ = strconv.ParseFloat(line[5], 64)
	}
	return p
}

// ReadPrices - read a csv quote file into Prices structure
func ReadPrices(filename string) Prices {
	csv, err := ioutil.ReadFile(filename)
	check(err)
	p := PricesFromCSV(string(csv))
	return p
}

// ToJSON - convert Prices to json string
func (p *Prices) ToJSON(indent bool) string {
	var j []byte
	if indent {
		j, _ = json.MarshalIndent(p, "", "  ")
	} else {
		j, _ = json.Marshal(p)
	}
	return string(j)
}

// WriteJSON - write Prices struct to json file
func (p *Prices) WriteJSON(filename string, indent bool) {
	json := p.ToJSON(indent)
	ba := []byte(json)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// ToCSV - convert Prices structure to csv string
func (s Symbols) ToCSV() string {

	var buffer bytes.Buffer

	buffer.WriteString("symbol,date,open,high,low,close,volume\n")

	for sym := 0; sym < len(s); sym++ {
		p := s[sym]
		for bar := range p.Close {
			str := fmt.Sprintf("%s,%s,%.2f,%.2f,%.2f,%.2f,%.0f\n",
				p.Symbol, p.Date[bar].Format("2006-01-02"), p.Open[bar], p.High[bar], p.Low[bar], p.Close[bar], p.Volume[bar])
			buffer.WriteString(str)
		}
	}

	return buffer.String()
}

// WriteCSV - save Symbols structure to file
func (s Symbols) WriteCSV(filename string) {
	csv := s.ToCSV()
	ba := []byte(csv)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// SymbolsFromCSV - parse csv quote string into Symbols array
func SymbolsFromCSV(csv string) Symbols {

	symbols := Symbols{}

	tmp := strings.Split(csv, "\n")
	numrows := len(tmp) - 1

	var index = make(map[string]int)
	for idx := 1; idx < numrows; idx++ {
		sym := strings.Split(tmp[idx], ",")[0]
		index[sym]++
	}

	row := 1
	for sym, len := range index {
		p := Prices{}
		p.Symbol = sym
		fmt.Printf("processing symbol %s, len=%d\n", sym, len)
		p.Date = make([]time.Time, len)
		p.Open = make([]float64, len)
		p.High = make([]float64, len)
		p.Low = make([]float64, len)
		p.Close = make([]float64, len)
		p.Volume = make([]float64, len)
		for bar := 0; bar < len; bar++ {
			line := strings.Split(tmp[row], ",")
			fmt.Println(line)
			p.Date[bar], _ = time.Parse("2006-01-02", line[1])
			p.Open[bar], _ = strconv.ParseFloat(line[2], 64)
			p.High[bar], _ = strconv.ParseFloat(line[3], 64)
			p.Low[bar], _ = strconv.ParseFloat(line[4], 64)
			p.Close[bar], _ = strconv.ParseFloat(line[5], 64)
			p.Volume[bar], _ = strconv.ParseFloat(line[6], 64)
			row++
		}
		fmt.Println(p)
		symbols = append(symbols, p)
	}
	return symbols
}

// ReadSymbols - read a csv quote file into Symbols array
func ReadSymbols(filename string) Symbols {
	csv, err := ioutil.ReadFile(filename)
	check(err)
	s := SymbolsFromCSV(string(csv))
	return s
}

// NewYahoo - Yahoo historical prices for a symbol
func NewYahoo(symbol, startDate, endDate string, period Period, adjustPrice bool) (*Prices, error) {

	from := parseDTString(startDate)

	var to time.Time
	if endDate == "" {
		to = time.Now()
	} else {
		to = parseDTString(endDate)
	}

	prices := &Prices{Symbol: symbol}

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
		prices.Date[bar] = d
		prices.Open[bar] = o * factor
		prices.High[bar] = h * factor
		prices.Low[bar] = l * factor
		prices.Close[bar] = c * factor
		prices.Volume[bar] = v
	}

	return prices, nil
}

// NewYahooSymbols - create a list of prices from symbols in file
func NewYahooSymbols(filename, startDate, endDate string, period Period, adjustPrice bool) (Symbols, error) {

	symbols := Symbols{}
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		sym := scanner.Text()
		fmt.Println("sym=" + sym)
		p, _ := NewYahoo(sym, startDate, endDate, period, adjustPrice)
		symbols = append(symbols, *p)
	}
	return symbols, nil
}
