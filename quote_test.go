package quote

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// assert fails the test if the condition is false.
func assert(t *testing.T, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(t *testing.T, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n", filepath.Base(file), line, err.Error())
		t.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\texp: %#v\n\tgot: %#v\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

// Get stock price history.
func TestStockPriceHistory(t *testing.T) {
	symbol := "AAPL"
	from, _ := time.Parse("2006-01-02", "2015-01-02")
	to, _ := time.Parse("2006-01-02", "2015-02-13")
	prices, err := GetDailyHistory(symbol, from, to, true)

	ok(t, err)
	equals(t, int(1), int(prices.Date[0].Month()))
	equals(t, int(2), int(prices.Date[0].Day()))
	equals(t, int(2015), int(prices.Date[0].Year()))

	equals(t, int(2), int(prices.Date[len(prices.Date)-1].Month()))
	equals(t, int(13), int(prices.Date[len(prices.Date)-1].Day()))
	equals(t, int(2015), int(prices.Date[len(prices.Date)-1].Year()))

}
