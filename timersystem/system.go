package timersystem

import (
	"container/heap"
	"fmt"
	"time"
)

type Timers struct {
	cur *time.Timer
	sh  chan *scheduleHeap
}

func NewTimers() Timers {
	ts := Timers{}
	ts.sh = make(chan *scheduleHeap, 1)
	ts.sh <- &scheduleHeap{}
	return ts
}

func (t *Timers) SetTimer(waitTime time.Duration, callback func()) {
	if t.cur != nil {
		t.cur.Stop()
		t.cur = nil
	}

	t.cur = time.AfterFunc(waitTime, func() {
		t.runCallback(callback)
	})
}

func (t *Timers) runCallback(callback func()) {
	go callback()
	t.cur = nil
}

func (t *Timers) GetCurrentTime() time.Time {
	return time.Now()
}

func (t *Timers) AddTimer(waitTime time.Duration, callback func(t int64)) {
	triggerTime := t.GetCurrentTime().Add(waitTime)
	sh := <-t.sh
	heap.Push(sh, &schedule{t: triggerTime, callback: callback})
	t.sh <- sh
	t.scheduleTimer()
}

func (t *Timers) scheduleTimer() {
	sh := <-t.sh
	if sh.Len() > 0 {
		first := sh.Peek()
		runFunc := func(in *schedule) {
			sh := <-t.sh
			if sh.Len() > 0 {
				check := sh.Peek()
				if check == in {
					heap.Pop(sh)
					check.callback(check.t.UnixNano())
				} else {
					fmt.Printf("skipping %v %v\n", check.t.UnixNano(), in.t.UnixNano())
				}
			}
			t.sh <- sh
			t.scheduleTimer()
		}
		current := t.GetCurrentTime()
		if first.t.Unix() < current.Unix() {
			t.SetTimer(0, func() {
				runFunc(first)
			})
		} else {
			t.SetTimer(first.t.Sub(current), func() {
				runFunc(first)
			})
		}
	}
	t.sh <- sh
}

type schedule struct {
	t        time.Time
	callback func(int64)
}

type scheduleHeap struct {
	buf []*schedule
}

func (sh scheduleHeap) Len() int {
	return len(sh.buf)
}

func (sh scheduleHeap) Less(i, j int) bool {
	return sh.buf[i].t.UnixNano() < sh.buf[j].t.UnixNano()
}

func (sh scheduleHeap) Swap(i, j int) {
	sh.buf[i], sh.buf[j] = sh.buf[j], sh.buf[i]
}

func (sh *scheduleHeap) Push(it interface{}) {
	sh.buf = append(sh.buf, it.(*schedule))
}

func (sh *scheduleHeap) Pop() interface{} {
	it := sh.buf[len(sh.buf)-1]
	sh.buf = sh.buf[:len(sh.buf)-1]
	return it
}

func (sh scheduleHeap) Peek() *schedule {
	return sh.buf[0]
}
