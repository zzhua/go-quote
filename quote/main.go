/*
Package quote is free quote downloader library and cli

Downloads daily/weekly/monthly/yearly historical price quotes from Yahoo
and daily/intraday data from Google

Copyright 2016 Mark Chenoweth
Licensed under terms of MIT license (see LICENSE)

Usage:
  quote -h | -help
  quote -v | -version
  quote etf|nyse|amex|nasdaq [-output=<outputFile>]
  quote [-years=<years>|(-start=<startDate> [-end=<endDate>])] [flags] [-infile=<inputFile>|<symbol>...]

Options:
  -h -help              show help
  -v -version           show version
  -years=<years>        number of years to download [default=5]
  -start=<startDate>    yyyy[-[mm-[dd]]]
  -end=<endDate>        yyyy[-[mm-[dd]]]
  -infile=<inputFile>   list of symbols to download
  -outfile=<outputFile> output filename
  -period=<period>      1m|5m|15m|30m|1h|d|w|m [default=d]
  -source=<source>      yahoo|google [default=yahoo]
  -format=<outFormat>   (csv|json) [default: json]
  -adjust=<bool>        adjust yahoo prices
  -all=<bool>           all in one file (true|false) [default: true]
*/
package main

// TODO:
// testing
// pacing flag
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

const (
	version    = "0.1"
	dateFormat = "2006-01-02"
)

var yearsFlag int
var startFlag string
var endFlag string
var periodFlag string
var sourceFlag string
var infileFlag string
var outfileFlag string
var formatFlag string
var allFlag bool
var adjustFlag bool
var versionFlag bool

func init() {
	flag.IntVar(&yearsFlag, "years", 5, "number of years to download")
	flag.StringVar(&startFlag, "begin", "", "start date (yyyy[-mm[-dd]])")
	flag.StringVar(&endFlag, "end", "", "end date (yyyy[-mm[-dd]])")
	flag.StringVar(&periodFlag, "period", "d", "1m|5m|15m|30m|1h|d|w|m")
	flag.StringVar(&sourceFlag, "source", "yahoo", "yahoo|google")
	flag.StringVar(&infileFlag, "infile", "", "input filename")
	flag.StringVar(&outfileFlag, "outfile", "", "output filename")
	flag.StringVar(&formatFlag, "format", "csv", "csv|json")
	flag.BoolVar(&allFlag, "all", true, "all output in one file")
	flag.BoolVar(&adjustFlag, "v", true, "adjust Yahoo prices")
	flag.BoolVar(&versionFlag, "v", false, "show version")
	flag.BoolVar(&versionFlag, "version", false, "show version")
	flag.Parse()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	// determine symbols
	var symbols []string
	if infileFlag != "" {
		raw, err := ioutil.ReadFile(infileFlag)
		check(err)
		symbols = strings.Split(string(raw), "\n")
	} else {
		symbols = flag.Args()
	}
	if len(symbols) == 0 {
		panic(fmt.Errorf("no symbols"))
	}

	// handle exchanges
	switch symbols[0] {
	case "etf":
		quote.NewEtfFile(outfileFlag)
		os.Exit(0)
	case "nyse":
		quote.NewExchangeFile("nyse", outfileFlag)
		os.Exit(0)
	case "nasdaq":
		quote.NewExchangeFile("nasdaq", outfileFlag)
		os.Exit(0)
	case "amex":
		quote.NewExchangeFile("amex", outfileFlag)
		os.Exit(0)
	}

	// validate source
	if sourceFlag != "yahoo" && sourceFlag != "google" {
		panic(fmt.Errorf("invalid source"))
	}

	// validate period
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

	// determine begin/end times
	var from, to time.Time
	if startFlag != "" {
		from = quote.ParseDTString(startFlag)
		if endFlag != "" {
			to = quote.ParseDTString(endFlag)
		} else {
			to = time.Now()
		}
	} else {
		to = time.Now()
		from = to.Add(-time.Duration(int(time.Hour) * 24 * 365 * yearsFlag))
	}

	// main output
	if allFlag {
		// output all in one file
		quotes := quote.Quotes{}
		if sourceFlag == "yahoo" {
			quotes, _ = quote.NewQuotesFromYahooSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period, adjustFlag)
		} else if sourceFlag == "google" {
			quotes, _ = quote.NewQuotesFromGoogleSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period)
		}
		if formatFlag == "csv" {
			quotes.WriteCSV(outfileFlag)
		} else if formatFlag == "json" {
			quotes.WriteJSON(outfileFlag, false)
		}
	} else {
		// output individual symbol files
		for _, sym := range symbols {
			var q quote.Quote
			if sourceFlag == "yahoo" {
				q, _ = quote.NewQuoteFromYahoo(sym, from.Format(dateFormat), to.Format(dateFormat), period, adjustFlag)
			} else if sourceFlag == "google" {
				q, _ = quote.NewQuoteFromGoogle(sym, from.Format(dateFormat), to.Format(dateFormat), period)
			}
			if formatFlag == "csv" {
				q.WriteCSV(outfileFlag)
			} else if formatFlag == "json" {
				q.WriteJSON(outfileFlag, false)
			}
		}
	}
}
