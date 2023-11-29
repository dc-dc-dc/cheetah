package market

import "errors"

type ContextKey string

const (
	ContextKeySymbol ContextKey = "market:symbol"
)

var (
	ErrSymbolNotFound = errors.New("symbol not found")
)

type Interval int

const (
	Interval1Minute Interval = iota
	Interval2Minute
	Interval5Minute
	Interval15Minute
	Interval30Minute
	Interval60Minute
	Interval90Minute
	Interval1Hour
	Interval1Day
	Interval5Day
	Interval1Week
	Interval1Month
	Interval3Month
)

func (i Interval) String() string {
	switch i {
	case Interval1Minute:
		return "1m"
	case Interval2Minute:
		return "2m"
	case Interval5Minute:
		return "5m"
	case Interval15Minute:
		return "15m"
	case Interval30Minute:
		return "30m"
	case Interval60Minute:
		return "60m"
	case Interval90Minute:
		return "90m"
	case Interval1Hour:
		return "1h"
	case Interval1Day:
		return "1d"
	case Interval5Day:
		return "5d"
	case Interval1Week:
		return "1wk"
	case Interval1Month:
		return "1mo"
	case Interval3Month:
		return "3mo"
	}
	return ""
}
