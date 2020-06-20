package timersystem

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimers(t *testing.T) {
	ts := NewTimers()
	var resArr []int64
	res := make(chan []int64, 1)
	res <- resArr
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		go ts.AddTimer(0, func(t int64) {
			r := <-res
			r = append(r, t)
			res <- r
		})
	}

	time.Sleep(time.Second)

	for {
		sh := <-ts.sh
		fmt.Printf("len:%d\n", sh.Len())
		if sh.Len() == 0 {
			break
		}
		ts.sh <- sh
		time.Sleep(time.Second)
	}
	r := <-res
	fmt.Printf("%v", r)
	rc := make([]int64, len(r))
	copy(rc, r)
	sort.Slice(rc, func(i, j int) bool {
		return rc[i] < rc[j]
	})
	assert.Equal(t, rc, r)
}
