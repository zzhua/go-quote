package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/markcheno/go-quote"
	//"time"
)

const version = "0.1"

func main() {
	usage := `
free quote downloader

Usage:
  quote -h | --help
  quote -v | --version
  quote <symbol>... [-y years|(-b <beginDate> [-e <endDate>])] [-ofldsap]
  quote -i <symFile> [-y <years>|(-b <beginDate> [-e <endDate>])] [-ofldsap]
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
  -d --date-format <dateFormat>  [default: yyyy-mm-dd]
  -a --all-in-one <allInOne>     All in one file (true|false) [default: false]`

	arguments, _ := docopt.Parse(usage, nil, true, version, false)
	fmt.Println(arguments)

	//first := true
	for _, sym := range arguments["<symbol>"].([]string) {
		q, _ := quote.NewYahooYears(sym, 1, quote.Daily, true)
		//fmt.Print(q.ToCSV(first, true))
		fmt.Print(q.ToJSON(true))
		//sfirst = false
	}
}
