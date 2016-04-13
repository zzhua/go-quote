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
  quote [-years=<years>|(-start=<datestr> [-end=<datestr>])] [options] [-infile=<filename>|<symbol> ...]

Options:
  -h -help             show help
  -v -version          show version
  -years=<years>       number of years to download [default=5]
  -start=<datestr>     yyyy[-[mm-[dd]]]
  -end=<datestr>       yyyy[-[mm-[dd]]] [default=today]
  -infile=<filename>   list of symbols to download
  -outfile=<filename>  output filename
  -period=<period>     1m|5m|15m|30m|1h|d|w|m [default=d]
  -source=<source>     yahoo|google [default=yahoo]
  -format=<format>     (csv|json) [default=csv]
  -adjust=<bool>       adjust yahoo prices [default=true]
  -all=<bool>          all in one file (true|false) [default=false]
  -log=<dest>          filename|stdout|stderr|discard [default=stdout]
  -delay=<ms>          delay in milliseconds between quote requests
*/
package main

import (
	"flag"
	"fmt"
	"github.com/markcheno/go-quote"
	"io/ioutil"
	"os"
	"time"
)

const (
	version    = "0.1"
	dateFormat = "2006-01-02"
)

type quoteflags struct {
	years   int
	delay   int
	start   string
	end     string
	period  string
	source  string
	infile  string
	outfile string
	format  string
	log     string
	all     bool
	adjust  bool
	version bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkFlags(flags quoteflags) error {

	// validate source
	if flags.source != "yahoo" && flags.source != "google" {
		return fmt.Errorf("invalid source, must be either 'yahoo' or 'google'")
	}

	// validate period
	if flags.source == "yahoo" &&
		(flags.period == "1m" || flags.period == "5m" || flags.period == "15m" || flags.period == "30m" || flags.period == "60m") {
		return fmt.Errorf("invalid source for yahoo, must be 'd' or 'w' or 'm'")

	}
	if flags.source == "google" && (flags.period == "w" || flags.period == "m") {
		return fmt.Errorf("invalid source for google, must be '1m' or '5m' or '15m' or '30m' or '60m' or 'd'")

	}
	return nil
}

func setOutput(flags quoteflags) error {
	var err error
	if flags.log == "stdout" {
		quote.Log.SetOutput(os.Stdout)
	} else if flags.log == "stderr" {
		quote.Log.SetOutput(os.Stderr)
	} else if flags.log == "discard" {
		quote.Log.SetOutput(ioutil.Discard)
	} else {
		var f *os.File
		f, err = os.OpenFile(flags.log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		quote.Log.SetOutput(f)
	}
	return err
}

func getSymbols(flags quoteflags, args []string) ([]string, error) {

	var err error
	var symbols []string

	if flags.infile != "" {
		symbols, err = quote.NewSymbolsFromFile(flags.infile)
		if err != nil {
			return symbols, err
		}
	} else {
		symbols = args
	}

	// make sure we found some symbols
	if len(symbols) == 0 {
		return symbols, fmt.Errorf("no symbols specified")
	}

	// validate outfileFlag
	if len(symbols) > 1 && flags.outfile != "" && !flags.all {
		return symbols, fmt.Errorf("outfile not valid with multiple symbols\nuse -all=true")
	}

	return symbols, nil
}

func getPeriod(periodFlag string) quote.Period {
	period := quote.Daily
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
	}
	return period
}

func getTimes(flags quoteflags) (time.Time, time.Time) {
	// determine start/end times
	to := quote.ParseDateString(flags.end)
	var from time.Time
	if flags.start != "" {
		from = quote.ParseDateString(flags.start)
	} else { // use years
		from = to.Add(-time.Duration(int(time.Hour) * 24 * 365 * flags.years))
	}
	return from, to
}

func outputAll(symbols []string, flags quoteflags) error {
	// output all in one file
	from, to := getTimes(flags)
	period := getPeriod(flags.period)
	quotes := quote.Quotes{}
	var err error
	if flags.source == "yahoo" {
		quotes, err = quote.NewQuotesFromYahooSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period, flags.adjust)
	} else if flags.source == "google" {
		quotes, err = quote.NewQuotesFromGoogleSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period)
	}
	if err != nil {
		return err
	}

	if flags.format == "csv" {
		err = quotes.WriteCSV(flags.outfile)
	} else if flags.format == "json" {
		err = quotes.WriteJSON(flags.outfile, false)
	}
	return err
}

func outputIndividual(symbols []string, flags quoteflags) error {
	// output individual symbol files

	from, to := getTimes(flags)
	period := getPeriod(flags.period)

	for _, sym := range symbols {
		var q quote.Quote
		if flags.source == "yahoo" {
			q, _ = quote.NewQuoteFromYahoo(sym, from.Format(dateFormat), to.Format(dateFormat), period, flags.adjust)
		} else if flags.source == "google" {
			q, _ = quote.NewQuoteFromGoogle(sym, from.Format(dateFormat), to.Format(dateFormat), period)
		}
		if flags.format == "csv" {
			_ = q.WriteCSV(flags.outfile)
		} else if flags.format == "json" {
			_ = q.WriteJSON(flags.outfile, false)
		}
		time.Sleep(quote.Delay * time.Millisecond)
	}
	return nil
}

func handleCommand(cmd string, flags quoteflags) bool {
	// handle exchange special commands
	handled := false
	switch cmd {
	case "etf":
		quote.NewEtfFile(flags.outfile)
		handled = true
	case "nyse":
		quote.NewExchangeFile("nyse", flags.outfile)
		handled = true
	case "nasdaq":
		quote.NewExchangeFile("nasdaq", flags.outfile)
		handled = true
	case "amex":
		quote.NewExchangeFile("amex", flags.outfile)
		handled = true
	}
	return handled
}

func main() {

	var err error
	var symbols []string
	var flags quoteflags

	flag.IntVar(&flags.years, "years", 5, "number of years to download")
	flag.IntVar(&flags.delay, "delay", 100, "milliseconds to delay between requests")
	flag.StringVar(&flags.start, "start", "", "start date (yyyy[-mm[-dd]])")
	flag.StringVar(&flags.end, "end", "", "end date (yyyy[-mm[-dd]])")
	flag.StringVar(&flags.period, "period", "d", "1m|5m|15m|30m|1h|d|w|m")
	flag.StringVar(&flags.source, "source", "yahoo", "yahoo|google")
	flag.StringVar(&flags.infile, "infile", "", "input filename")
	flag.StringVar(&flags.outfile, "outfile", "", "output filename")
	flag.StringVar(&flags.format, "format", "csv", "csv|json")
	flag.StringVar(&flags.log, "log", "stdout", "<filename>|stdout")
	flag.BoolVar(&flags.all, "all", false, "all output in one file")
	flag.BoolVar(&flags.adjust, "adjust", true, "adjust Yahoo prices")
	flag.BoolVar(&flags.version, "v", false, "show version")
	flag.BoolVar(&flags.version, "version", false, "show version")
	flag.Parse()

	if flags.version {
		fmt.Println(version)
		os.Exit(0)
	}

	quote.Delay = time.Duration(flags.delay)

	err = setOutput(flags)
	check(err)

	err = checkFlags(flags)
	check(err)

	symbols, err = getSymbols(flags, flag.Args())
	check(err)

	// check for and handled special commands
	if handleCommand(symbols[0], flags) {
		os.Exit(0)
	}

	// main output
	if flags.all {
		err = outputAll(symbols, flags)
	} else {
		err = outputIndividual(symbols, flags)
	}
}
