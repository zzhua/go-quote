# go-quote

[![GoDoc](http://godoc.org/github.com/markcheno/go-quote?status.svg)](http://godoc.org/github.com/markcheno/go-quote) 


A free quote downloader library and cli

Downloads daily/weekly/monthly/yearly historical price quotes from Yahoo
and daily/intraday data from Google

Copyright 2016 Mark Chenoweth
Licensed under terms of MIT license (see LICENSE)

Install CLI utility (quote) with:

```bash
go install github.com/markcheno/go-quote/quote
```

```bash
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
```

## CLI Examples

```bash
# display usage
quote -help

# downloads 5 years of Yahoo SPY history to spy.csv 
quote spy

# downloads 1 year of Yahoo SPY & AAPL history to quotes.csv 
quote -years=1 -outfile=quotes.csv spy aapl

# downloads 2 years of Google SPY & AAPL history to spy.csv and aapl.csv 
quote -years=2 -all=false -source=google spy aapl

# downloads full etf symbol list to etf.txt, also works for nasdaq,nyse,amex
quote etf

# downloads quote history for symbols in myquotelist.txt from 2000 to today into myquotes.csv
quote -start=2000 -infile=myquotelist.txt -outfile=myquotes.csv 
```

## Install library

Install the package with:

```bash
go get github.com/markcheno/go-quote
```

## Library example

```go
import "github.com/markcheno/go-quote"

// downloads weekly, adjusted SPY history from 2000 to 2010
spy := quote.NewQuoteYahoo("spy","2000","2010",quote.Weekly,true)

// downloads 5 minute GLD history for the past 10 days
gld := quote.NewQuoteGoogle("gld","","",quote.Min5)
```

## License

MIT License  - see LICENSE for more details
