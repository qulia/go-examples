package rangemodule

import "github.com/qulia/go-qulia/lib/tree"

type RangeModule struct {
	st tree.SegmentTreeInterface
}

func Constructor() RangeModule {
	rm := RangeModule{}
	rm.st = tree.NewSegmentTree(trackingQueryEvalFunc, func() interface{} { return true })
	rm.RemoveRange(0, 1e9+1)
	return rm
}

func (rm *RangeModule) AddRange(left int, right int) {
	rm.st.UpdateRange(left, right-1, func(_ interface{}) interface{} {
		return true
	})
}

func (rm *RangeModule) QueryRange(left int, right int) bool {
	return rm.st.QueryRange(left, right-1).(bool)
}

func (this *RangeModule) RemoveRange(left int, right int) {
	this.st.UpdateRange(left, right-1, func(_ interface{}) interface{} {
		return false
	})
}

func trackingQueryEvalFunc(a, b interface{}) interface{} {
	if a == nil || b == nil {
		return false
	}

	return a.(bool) && b.(bool)
}
