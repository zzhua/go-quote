package quote

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	// Yearly time period
	Yearly Period = "y"
)

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

// CSV - convert Quote structure to csv string
func (q *Quote) CSV() string {

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
func (q *Quote) WriteCSV(filename string) {
	csv := q.CSV()
	ba := []byte(csv)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// NewQuoteFromCSV - parse csv quote string into Quote structure
func NewQuoteFromCSV(csv string) Quote {

	q := Quote{}
	tmp := strings.Split(csv, "\n")
	numrows := len(tmp) - 1
	q.Date = make([]time.Time, numrows-1)
	q.Open = make([]float64, numrows-1)
	q.High = make([]float64, numrows-1)
	q.Low = make([]float64, numrows-1)
	q.Close = make([]float64, numrows-1)
	q.Volume = make([]float64, numrows-1)

	for row, bar := 1, 0; row < numrows; row, bar = row+1, bar+1 {
		line := strings.Split(tmp[row], ",")
		q.Date[bar], _ = time.Parse("2006-01-02 15:04", line[0])
		q.Open[bar], _ = strconv.ParseFloat(line[1], 64)
		q.High[bar], _ = strconv.ParseFloat(line[2], 64)
		q.Low[bar], _ = strconv.ParseFloat(line[3], 64)
		q.Close[bar], _ = strconv.ParseFloat(line[4], 64)
		q.Volume[bar], _ = strconv.ParseFloat(line[5], 64)
	}
	return q
}

// NewQuoteFromCSVFile - parse csv quote file into Quote structure
func NewQuoteFromCSVFile(filename string) Quote {
	csv, err := ioutil.ReadFile(filename)
	check(err)
	return NewQuoteFromCSV(string(csv))

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
func (q Quote) WriteJSON(filename string, indent bool) {
	json := q.JSON(indent)
	ba := []byte(json)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// NewQuoteFromJSON - parse json quote string into Quote structure
func NewQuoteFromJSON(jsn string) Quote {
	q := Quote{}
	err := json.Unmarshal([]byte(jsn), &q)
	check(err)
	return q
}

// NewQuoteFromJSONFile - parse json quote string into Quote structure
func NewQuoteFromJSONFile(filename string) Quote {
	jsn, err := ioutil.ReadFile(filename)
	check(err)
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
func (q Quotes) WriteCSV(filename string) {
	csv := q.CSV()
	ba := []byte(csv)
	err := ioutil.WriteFile(filename, ba, 0644)
	check(err)
}

// NewQuotesFromCSV - parse csv quote string into Quotes array
func NewQuotesFromCSV(csv string) Quotes {

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
		q := Quote{}
		q.Symbol = sym
		q.Date = make([]time.Time, len)
		q.Open = make([]float64, len)
		q.High = make([]float64, len)
		q.Low = make([]float64, len)
		q.Close = make([]float64, len)
		q.Volume = make([]float64, len)
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
	return quotes
}

// NewQuotesFromCSVFile - parse csv quote file into Quotes array
func NewQuotesFromCSVFile(filename string) Quotes {
	csv, err := ioutil.ReadFile(filename)
	check(err)
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
func (q Quotes) WriteJSON(filename string, indent bool) {
	jsn := q.JSON(indent)
	err := ioutil.WriteFile(filename, []byte(jsn), 0644)
	check(err)
}

// NewQuotesFromJSON - parse json quote string into Quote structure
func NewQuotesFromJSON(jsn string) Quotes {
	quotes := Quotes{}
	err := json.Unmarshal([]byte(jsn), &quotes)
	check(err)
	return quotes
}

// NewQuotesFromJSONFile - parse json quote string into Quote structure
func NewQuotesFromJSONFile(filename string) Quotes {
	jsn, err := ioutil.ReadFile(filename)
	check(err)
	return NewQuotesFromJSON(string(jsn))
}

// NewQuoteFromYahoo - Yahoo historical prices for a symbol
func NewQuoteFromYahoo(symbol, startDate, endDate string, period Period, adjustQuote bool) (Quote, error) {

	from := parseDTString(startDate)

	var to time.Time
	if endDate == "" {
		to = time.Now()
	} else {
		to = parseDTString(endDate)
	}

	quote := Quote{Symbol: symbol}

	url := fmt.Sprintf(
		"http://ichart.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%d&e=%d&f=%d&g=%s&ignore=.csv",
		symbol,
		from.Month()-1, from.Day(), from.Year(),
		to.Month()-1, to.Day(), to.Year(),
		period)

	resp, err := http.Get(url)
	if err != nil {
		return quote, err
	}
	defer resp.Body.Close()

	var csvdata [][]string
	reader := csv.NewReader(resp.Body)
	csvdata, err = reader.ReadAll()
	if err != nil {
		return quote, err
	}

	numrows := len(csvdata) - 1
	quote.Date = make([]time.Time, numrows)
	quote.Open = make([]float64, numrows)
	quote.High = make([]float64, numrows)
	quote.Low = make([]float64, numrows)
	quote.Close = make([]float64, numrows)
	quote.Volume = make([]float64, numrows)

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
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		sym := scanner.Text()
		quote, _ := NewQuoteFromYahoo(sym, startDate, endDate, period, adjustQuote)
		quotes = append(quotes, quote)
	}
	return quotes, nil
}

// NewQuoteFromGoogle - Google daily/intraday historical prices for a symbol
func NewQuoteFromGoogle(symbol, startDate, endDate string, period Period) (Quote, error) {

	from := parseDTString(startDate)

	var to time.Time
	if endDate == "" {
		to = time.Now()
	} else {
		to = parseDTString(endDate)
	}

	quote := Quote{Symbol: symbol}

	var args string
	if period == Daily {
		args = fmt.Sprintf(
			"http://www.google.com/finance/historical?q=%s&startdate=%s&enddate=%s&output=csv",
			symbol,
			url.QueryEscape(from.Format("Jan 2, 2006")),
			url.QueryEscape(to.Format("Jan 2, 2006")))

	} else if period == Min5 || period == Min15 || period == Min30 || period == Min60 {

		args = fmt.Sprintf(
			"http://www.google.com/finance/getprices?q=%s&i=%s&p=10d&f=d,o,h,l,c,v",
			strings.ToUpper(symbol),
			period)

	} else {
		return quote, fmt.Errorf("invalid period")
	}

	fmt.Println(args)
	resp, err := http.Get(args)
	if err != nil {
		return quote, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)

	var csvdata [][]string

	var tmp string
	if period == Daily {
		tmp = strings.Join(strings.Split(string(contents), "\n")[1:], "\n")
	} else {
		tmp = strings.Join(strings.Split(string(contents), "\n")[7:], "\n")
	}

	reader := csv.NewReader(strings.NewReader(tmp))
	csvdata, err = reader.ReadAll()
	if err != nil {
		return quote, err
	}

	numrows := len(csvdata)
	quote.Date = make([]time.Time, numrows)
	quote.Open = make([]float64, numrows)
	quote.High = make([]float64, numrows)
	quote.Low = make([]float64, numrows)
	quote.Close = make([]float64, numrows)
	quote.Volume = make([]float64, numrows)

	fmt.Printf("numrows = %d, len(csvdata)= % d\n", numrows, len(csvdata))

	var day int64

	for row := 0; row < numrows; row++ {

		var d time.Time
		var o, h, l, c, v float64

		if period == Daily {
			d, _ = time.Parse("2-Jan-06", csvdata[row][0])
			o, _ = strconv.ParseFloat(csvdata[row][1], 64)
			h, _ = strconv.ParseFloat(csvdata[row][2], 64)
			l, _ = strconv.ParseFloat(csvdata[row][3], 64)
			c, _ = strconv.ParseFloat(csvdata[row][4], 64)
			v, _ = strconv.ParseFloat(csvdata[row][5], 64)

		} else {
			fmt.Println(csvdata[row])
			c, _ = strconv.ParseFloat(csvdata[row][1], 64)
			h, _ = strconv.ParseFloat(csvdata[row][2], 64)
			l, _ = strconv.ParseFloat(csvdata[row][3], 64)
			o, _ = strconv.ParseFloat(csvdata[row][4], 64)
			v, _ = strconv.ParseFloat(csvdata[row][5], 64)

			var offset int64
			z := csvdata[row][0]
			if z[0] == 'a' {
				day, _ = strconv.ParseInt(z[1:], 10, 64)
			} else {
				offset, _ = strconv.ParseInt(z, 10, 64)
			}
			seconds, _ := strconv.ParseInt(string(period), 10, 64)
			d = time.Unix(day+(seconds*offset), 0)

		}
		var bar int
		if period == Daily {
			bar = numrows - 1 - row // reverse the order
		} else {
			bar = row
		}
		quote.Date[bar] = d
		quote.Open[bar] = o
		quote.High[bar] = h
		quote.Low[bar] = l
		quote.Close[bar] = c
		quote.Volume[bar] = v
	}

	return quote, nil
}

// NewQuotesFromGoogle - create a list of prices from symbols in file
func NewQuotesFromGoogle(filename, startDate, endDate string, period Period) (Quotes, error) {

	quotes := Quotes{}
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		sym := scanner.Text()
		quote, _ := NewQuoteFromGoogle(sym, startDate, endDate, period)
		quotes = append(quotes, quote)
	}
	return quotes, nil
}
