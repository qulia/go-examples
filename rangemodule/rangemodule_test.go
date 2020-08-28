package rangemodule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeModule(t *testing.T) {
	rm := Constructor()
	rm.AddRange(10, 20)
	rm.RemoveRange(14, 16)
	assert.Equal(t, true, rm.QueryRange(10, 14))
	assert.Equal(t, false, rm.QueryRange(13, 15))
	assert.Equal(t, true, rm.QueryRange(16, 17))
}
