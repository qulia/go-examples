package rangemodule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeModule(t *testing.T) {
	rm := NewRangeModule()
	rm.AddRange(10, 20)
	rm.RemoveRange(14, 16)
	assert.Equal(t, 1, rm.QueryRange(10, 14))
	assert.Equal(t, 0, rm.QueryRange(13, 15))
	assert.Equal(t, 1, rm.QueryRange(16, 17))
}
