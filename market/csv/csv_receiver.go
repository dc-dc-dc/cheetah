package csv

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/dc-dc-dc/cheetah/market"
	"github.com/shopspring/decimal"
)

func NewCsvReceiver(dst io.Writer) market.MarketReceiver {
	writer := NewCsvWriter(&dst)
	header := []string{}
	return market.NewFunctionalReceiver(func(ctx context.Context, line market.MarketLine) error {
		headerSet := len(header) != 0
		elements := make([]string, 0, 5)
		elements = append(elements, line.Start.Format("2006-01-02 15:04:05"))
		elements = append(elements, line.Open.String())
		elements = append(elements, line.High.String())
		elements = append(elements, line.Low.String())
		elements = append(elements, line.Close.String())
		elements = append(elements, strconv.FormatInt(line.Volume, 10))
		cache, err := market.GetCache(ctx)
		if !headerSet {
			header = append(header, []string{"Date", "Open", "High", "Low", "Close", "Volume"}...)
		}
		if err == nil {
			cache.Range(func(key, val interface{}) bool {
				if !headerSet {
					header = append(header, key.(string))
				}
				if v, ok := val.(interface{}); ok {
					// loop over these...
					fmt.Printf("v: %v\n", v)
				} else {
					elements = append(elements, GetString(val))
				}
				return true
			})
		}
		if !headerSet {
			if err := writer.Write(header); err != nil {
				return err
			}
		}
		return writer.Write(elements)
	})
}

func GetString(t interface{}) string {
	switch t := t.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(t)
	case decimal.Decimal:
		return t.String()
	default:
		return ""
	}
}
