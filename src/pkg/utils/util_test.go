package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementInSlice(t *testing.T) {
	a := []int{1, 2, 3}
	exist := ElementInSlice(a, 1)
	assert.True(t, exist)
}
