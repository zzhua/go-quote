package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/markcheno/go-quote"
	//"time"
	//"encoding/json"
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
  -s --quote-source <source>     (yahoo|google|quandl) [default: yahoo]
  -o --out-file <outFile>        Output filename
  -f --out-format <outFormat>    (csv|json) [default: json]
  -l --log-file <logFile>        Log file [default: getq.log]
  -d --date-format <dateFormat>  [default: yyyy-mm-dd]
  -a --all-in-one <allInOne>     All in one file (true|false) [default: false]`

	arguments, _ := docopt.Parse(usage, nil, true, version, false)
	fmt.Println(arguments)

	q, _ := quote.HistoryForYears("aapl", 1, quote.Weekly, true)

	//from, _ := time.Parse("2006-01-02", "2015-01-02")
	//to, _ := time.Parse("2006-01-02", "2015-02-13")
	//q, _ := yahoo.GetDailyHistory("spy", from, to, true)

	fmt.Println(q.Bar(0))
	fmt.Println(q.Bar(len(q.Close) - 1))

	//j, _ := json.Marshal(q)
	//fmt.Println(string(j))

}
