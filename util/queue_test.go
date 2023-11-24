package util_test

import (
	"testing"

	"github.com/dc-dc-dc/cheetah/util"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	queue := util.NewQueue[int]()
	assert.NotNil(t, queue)
	queue.Pop()
	queue.PopLeft()
	queue.Push(1)
	queue.Push(2)
	queue.Push(3)
	assert.Equal(t, queue.Count(), 3)
	assert.Equal(t, queue.First(), 1)
	assert.Equal(t, queue.Last(), 3)
	assert.Equal(t, queue.Elements(), []int{1, 2, 3})
	queue.Pop()
	assert.Equal(t, queue.Count(), 2)
	assert.Equal(t, queue.First(), 1)
	assert.Equal(t, queue.Last(), 2)
	queue.PopLeft()
	assert.Equal(t, queue.Count(), 1)
	assert.Equal(t, queue.First(), 2)
	assert.Equal(t, queue.Last(), 2)
	queue.Pop()
	assert.Equal(t, queue.Count(), 0)
	assert.Zero(t, queue.First())
	assert.Zero(t, queue.Last())
	queue.Push(1)
	assert.Equal(t, queue.Count(), 1)
	queue.PopLeft()
	assert.Equal(t, queue.Count(), 0)
	assert.Zero(t, queue.First())
	assert.Zero(t, queue.Last())

}
