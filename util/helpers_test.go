package util_test

import (
	"testing"

	"github.com/dc-dc-dc/cheetah/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMinDecimal(t *testing.T) {
	testCases := []struct {
		left   int32
		right  int32
		answer int32
	}{
		{1, 2, 1},
		{2, 1, 1},
		{1, 1, 1},
		{-1, 1, -1},
	}
	for _, testCase := range testCases {
		assert.Equal(t, util.MinDecimal(decimal.NewFromInt32(testCase.left), decimal.NewFromInt32(testCase.right)).Cmp(decimal.NewFromInt32(testCase.answer)), 0)
	}
}

func TestMaxDecimal(t *testing.T) {
	testCases := []struct {
		left   int32
		right  int32
		answer int32
	}{
		{1, 2, 2},
		{2, 1, 2},
		{1, 1, 1},
		{-1, 1, 1},
	}
	for _, testCase := range testCases {
		assert.Equal(t, util.MaxDecimal(decimal.NewFromInt32(testCase.left), decimal.NewFromInt32(testCase.right)).Cmp(decimal.NewFromInt32(testCase.answer)), 0)
	}
}
