package rangemodule

import "github.com/qulia/go-qulia/lib/tree"

type RangeModule struct {
	st tree.SegmentTree[int]
}

func NewRangeModule() RangeModule {
	rm := RangeModule{}
	rm.st = tree.NewSegmentTree(aggFunc, 1)
	return rm
}

func (rm *RangeModule) AddRange(left int, right int) {
	rm.st.UpdateRange(left, right-1, func(_ int) int {
		return 1
	})
}

func (rm *RangeModule) RemoveRange(left int, right int) {
	rm.st.UpdateRange(left, right-1, func(_ int) int {
		return 0
	})
}

func (rm *RangeModule) QueryRange(left int, right int) int {
	return rm.st.QueryRange(left, right-1)
}

func aggFunc(a, b int) int {
	return (a + b) % 2
}
