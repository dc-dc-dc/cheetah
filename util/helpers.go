package util

import "github.com/shopspring/decimal"

func MinDecimal(i, j decimal.Decimal) decimal.Decimal {
	if i.LessThan(j) {
		return i
	}
	return j
}

func MaxDecimal(i, j decimal.Decimal) decimal.Decimal {
	if i.GreaterThan(j) {
		return i
	}
	return j
}
