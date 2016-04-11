# go-quote

[![GoDoc](http://godoc.org/github.com/markcheno/go-quote?status.svg)](http://godoc.org/github.com/markcheno/go-quote) 


A free quote downloader library and cli

Downloads daily/weekly/monthly historical price quotes from Yahoo
and daily/intraday data from Google

Still very much in alpha mode. Expect bugs and API changes. Comments/suggestions/pull requests welcome!

Copyright 2016 Mark Chenoweth

Install CLI utility (quote) with:

```bash
go install github.com/markcheno/go-quote/quote
```

```
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
```

## CLI Examples

```bash
# display usage
quote -help

# downloads 5 years of Yahoo SPY history to spy.csv 
quote spy

# downloads 1 year of Yahoo SPY & AAPL history to quotes.csv 
quote -years=1 -all=true -outfile=quotes.csv spy aapl

# downloads 2 years of Google SPY & AAPL history to spy.csv and aapl.csv 
quote -years=2 -source=google spy aapl

# downloads full etf symbol list to etf.txt, also works for nasdaq,nyse,amex
quote etf

# dowload fresh etf list and 5 years of etf data all in one file
quote etf && quote -all=true -outfile=etf.csv -infile=etf.txt 

# downloads quote history for symbols in myquotelist.txt from 2000 to today all into myquotes.csv
quote -start=2000  -all=true -outfile=myquotes.csv -infile=myquotelist.txt 
```

## Install library

Install the package with:

```bash
go get github.com/markcheno/go-quote
```

## Library example

```go
package main

import (
	"fmt"
	"github.com/markcheno/go-quote"
)

func main() { 
  // download weekly, adjusted SPY history from 2000-01-01 to 2009-12-31
	spy, _ := quote.NewQuoteFromYahoo("spy", "2000", "2009-12-31", quote.Weekly, true)
  // print csv formated data
	fmt.Print(spy.CSV())
}
```

## License

MIT License  - see LICENSE for more details
