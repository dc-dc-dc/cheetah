package util_test

import (
	"testing"

	"github.com/dc-dc-dc/cheetah/util"
	"github.com/stretchr/testify/assert"
)

func TestLinkedList(t *testing.T) {
	node := util.NewLinkedListNode(1)
	assert.NotNil(t, node)
}
