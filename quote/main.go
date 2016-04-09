package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/markcheno/go-quote"
)

const version = "0.1"

func main() {
	usage := `
free quote downloader

Usage:
  quote -h | --help
  quote -v | --version
  quote <symbol>... [-y years|(-b <beginDate> [-e <endDate>])] [-oflsap]
  quote -i <symFile> [-y <years>|(-b <beginDate> [-e <endDate>])] [-oflsap]
  quote (etf|nyse|amex|nasdaq) [-ofl]

Options:
  -h --help                      Show help
  -v --version                   Show version
  -i --in-file <symbolFile>      List of symbols to download
  -y --years <years>             Number of years to download [default: 5]
  -b --begin-date <beginDate>    (yyyy-mm-dd)
  -e --end-date <endDate>        (yyyy-mm-dd|today) [default: today]
  -p --period <period>           (1m|5m|15m|30m|1h|d|w|m) [default: d]
  -s --quote-source <source>     (yahoo|google|quandl|brownian) [default: yahoo]
  -o --out-file <outFile>        Output filename
  -f --out-format <outFormat>    (csv|json) [default: json]
  -l --log-file <logFile>        Log file [default: quote.log]
  -a --all-in-one <allInOne>     All in one file (true|false) [default: false]`

	arguments, _ := docopt.Parse(usage, nil, true, version, false)
	fmt.Println(arguments)

	//for _, sym := range arguments["<symbol>"].([]string) {

	//q, _ := quote.NewYahoo(sym, "2015", "", quote.Daily, false)
	//q, _ := quote.NewYahooYears(sym, 1, quote.Daily, false)
	//csv := q.ToCSV()
	//fmt.Print(csv)

	//p := quote.ReadCSV(sym + ".csv")
	//fmt.Print(p.ToCSV(false, false))

	//}
	syms, _ := quote.NewYahooSymbols("list.txt", "2016-04", "", quote.Daily, false)
	fmt.Println(syms)

	csv := syms.ToCSV()
	fmt.Println(csv)

	syms2 := quote.SymbolsFromCSV(csv)
	fmt.Println(syms2)

	syms2.WriteCSV("symbols.csv")

	syms3 := quote.ReadSymbols("symbols.csv")

	fmt.Println(syms3)
}

// NewYahooYears - get Yahoo stock price history for a number of years
//func NewYahooYears(symbol string, years int, period Period, adjustPrice bool) (*Prices, error) {
//	to := time.Now()
//	from := to.Add(-time.Duration(int(time.Hour) * 24 * 365 * years))
//	layout := "2006-01-02 15:04:05"
//	return NewYahoo(symbol, from.Format(layout), to.Format(layout), period, adjustPrice)
//}
