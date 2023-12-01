package indicator_test

import (
	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
)

var (
	testingFileName = util.GetEnv("TESTING_DATA_FILE", "./testing_data.csv")
)

func IsInRange(test, cmp, delta decimal.Decimal) bool {
	return test.GreaterThanOrEqual(cmp.Sub(delta)) && test.LessThanOrEqual(cmp.Add(delta))
}
