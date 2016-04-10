package main

/*
free quote downloader

Usage:
  quote -h | -help
  quote -v | -version
  quote [-y <years>|(-b <beginDate> [-e <endDate>])] [flags] [-i <symFile>|etf|nyse|amex|nasdaq|<symbol>...]

Options:
  -h -help                 Show help
  -v -version              Show version
  -infile <symbolFile>  List of symbols to download
  -years <years>        Number of years to download [default: 5]
  -begin <beginDate>    yyyy-mm-dd
  -end <endDate>        yyyy-mm-dd
  -period <period>      1m|5m|15m|30m|1h|d|w|m [default: d]
  -source <source>      yahoo|google [default: yahoo]
  -outfile <outFile>    Output filename
  -format <outFormat>   (csv|json) [default: json]
  -all <allInOne>      All in one file (true|false) [default: true]`
*/
// TODO:
// version flag
// yahoo adjust prices flag, pacing flag
// stdout/stdin? piping
// log file

import (
	"flag"
	"fmt"
	"github.com/markcheno/go-quote"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const version = "0.1"
const dateFormat = "2006-01-02"

var yearsFlag int
var beginFlag string
var endFlag string
var periodFlag string
var sourceFlag string
var inFlag string
var outFlag string
var formatFlag string
var allFlag bool

func init() {
	const (
		yearsUsage  = "Number of years to download"
		beginUsage  = "Begin date (yyyy[-mm[-dd]])"
		endUsage    = "End date (yyyy[-mm[-dd]])"
		periodUsage = "1m|5m|15m|30m|1h|d|w|m"
		sourceUsage = "yahoo|google"
		inUsage     = "Input filename"
		outUsage    = "Output filename"
		formatUsage = "csv|json"
		allUsage    = "all output in one file"
	)
	//flag.IntVar(&yearsFlag, "y", 5, yearsUsage)
	flag.IntVar(&yearsFlag, "years", 5, yearsUsage)

	//flag.StringVar(&beginFlag, "b", "", beginUsage)
	flag.StringVar(&beginFlag, "begin", "", beginUsage)

	//flag.StringVar(&endFlag, "e", "", endUsage)
	flag.StringVar(&endFlag, "end", "", endUsage)

	//flag.StringVar(&periodFlag, "p", "d", periodUsage)
	flag.StringVar(&periodFlag, "period", "d", periodUsage)

	//flag.StringVar(&sourceFlag, "s", "yahoo", sourceUsage)
	flag.StringVar(&sourceFlag, "source", "yahoo", sourceUsage)

	//flag.StringVar(&inFlag, "i", "", inUsage)
	flag.StringVar(&inFlag, "infile", "", inUsage)

	//flag.StringVar(&outFlag, "o", "", outUsage)
	flag.StringVar(&outFlag, "outfile", "", outUsage)

	//flag.StringVar(&formatFlag, "f", "csv", formatUsage)
	flag.StringVar(&formatFlag, "format", "csv", formatUsage)

	//flag.BoolVar(&allFlag, "a", true, allUsage)
	flag.BoolVar(&allFlag, "all", true, allUsage)

	flag.Parse()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// determine symbols
	var symbols []string
	if inFlag != "" {
		raw, err := ioutil.ReadFile(inFlag)
		check(err)
		symbols = strings.Split(string(raw), "\n")
	} else {
		symbols = flag.Args()
	}
	if len(symbols) == 0 {
		panic(fmt.Errorf("no symbols"))
	}
	//fmt.Println(symbols)

	// determine outfile
	//fmt.Println("outFlag=" + outFlag)

	// validate source
	if sourceFlag != "yahoo" && sourceFlag != "google" {
		panic(fmt.Errorf("invalid source"))
	}

	// determine period
	var period quote.Period
	switch periodFlag {
	case "1m":
		period = quote.Min1
	case "5m":
		period = quote.Min5
	case "15m":
		period = quote.Min15
	case "30m":
		period = quote.Min30
	case "60m":
		period = quote.Min60
	case "d":
		period = quote.Daily
	case "w":
		period = quote.Weekly
	case "m":
		period = quote.Monthly
	case "y":
		period = quote.Yearly
	}
	//fmt.Println("period=" + period)

	// handle exchanges
	switch symbols[0] {
	case "etf":
		quote.NewEtfFile(outFlag)
		os.Exit(0)
	case "nyse":
		quote.NewExchangeFile("nyse", outFlag)
		os.Exit(0)
	case "nasdaq":
		quote.NewExchangeFile("nasdaq", outFlag)
		os.Exit(0)
	case "amex":
		quote.NewExchangeFile("amex", outFlag)
		os.Exit(0)
	}

	// determine begin/end times
	var from, to time.Time

	if beginFlag != "" {

		from = quote.ParseDTString(beginFlag)

		if endFlag != "" {
			to = quote.ParseDTString(endFlag)
		} else {
			to = time.Now()
		}
	} else {
		to = time.Now()
		from = to.Add(-time.Duration(int(time.Hour) * 24 * 365 * yearsFlag))
	}
	//fmt.Printf("from=%s, to=%s", from, to)

	if allFlag {
		quotes := quote.Quotes{}
		if sourceFlag == "yahoo" {
			quotes, _ = quote.NewQuotesFromYahooSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period, true)
		} else if sourceFlag == "google" {
			quotes, _ = quote.NewQuotesFromGoogleSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period)
		}
		if formatFlag == "csv" {
			quotes.WriteCSV(outFlag)
		} else if formatFlag == "json" {
			quotes.WriteJSON(outFlag, false)
		}
		os.Exit(0) // done
	} else {
		// output individual symbol files
		for _, sym := range symbols {
			var q quote.Quote
			if sourceFlag == "yahoo" {
				q, _ = quote.NewQuoteFromYahoo(sym, from.Format(dateFormat), to.Format(dateFormat), period, true)
			} else if sourceFlag == "google" {
				q, _ = quote.NewQuoteFromGoogle(sym, from.Format(dateFormat), to.Format(dateFormat), period)
			}
			if formatFlag == "csv" {
				q.WriteCSV(outFlag)
			} else if formatFlag == "json" {
				q.WriteJSON(outFlag, false)
			}
		}
	}
}
