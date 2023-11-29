package util_test

import (
	"testing"

	"github.com/dc-dc-dc/cheetah/util"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	set := util.NewSet[string]()
	assert.NotNil(t, set)

	set.Add("a")
	set.Add("b")
	set.Add("c")
	assert.True(t, set.Contains("a"))
	assert.True(t, set.Contains("b"))
	assert.True(t, set.Contains("c"))
	assert.False(t, set.Contains("d"))
	set.Remove("a")
	assert.False(t, set.Contains("a"))
}
