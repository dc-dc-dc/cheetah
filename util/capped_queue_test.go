package util_test

import (
	"testing"

	"github.com/dc-dc-dc/cheetah/util"
	"github.com/stretchr/testify/assert"
)

func TestCappedQueue(t *testing.T) {
	queue := util.NewCappedQueue(3)
	assert.NotNil(t, queue)
	assert.Equal(t, queue.Cap(), 3)
	queue.Push(1)
	assert.False(t, queue.Full())
	queue.Push(2)
	queue.Push(3)
	assert.True(t, queue.Full())
	assert.Equal(t, queue.Elements(), []interface{}{1, 2, 3})
	queue.Push(4)
	assert.Equal(t, queue.Elements(), []interface{}{2, 3, 4})
}
