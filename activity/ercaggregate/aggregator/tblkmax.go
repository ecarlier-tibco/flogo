package aggregator

import (
	"sync"
	"time"
)

type TimeBlockMaximum struct {
	windowSize   time.Duration
	values       []float64
	windowMtx    *sync.Mutex
	startMtx     *sync.RWMutex
	windowActive bool
}

func init() {
	RegisterFactory("timeblockmax", NewTimeBlockMaximum)
}

func (ta *TimeBlockMaximum) Add(value float64) (bool, float64) {

	ta.windowMtx.Lock()
	ta.values = append(ta.values, value)
	ta.windowMtx.Unlock()

	if ta.startWindow() {
		time.Sleep(ta.windowSize * time.Millisecond)
		return true, ta.max()
	} else {
		return false, 0
	}
}

func (ta *TimeBlockMaximum) max() float64 {

	ta.windowMtx.Lock()

	var max = float64(0)

	count := len(ta.values)

	for i := 0; i < count; i++ {
		if max == 0 || ta.values[i] > max {
			max = ta.values[i]
		}
	}

	ta.resetWindow()

	return max
}

func (ta *TimeBlockMaximum) startWindow() bool {

	ta.startMtx.RLock()

	if ta.windowActive {
		ta.startMtx.RUnlock()
		return false
	}
	ta.startMtx.RUnlock()

	ta.startMtx.Lock()
	defer ta.startMtx.Unlock()

	if !ta.windowActive {
		ta.windowActive = true
		return true
	}

	return false
}

func (ta *TimeBlockMaximum) resetWindow() {
	ta.values = nil
	ta.windowMtx.Unlock()

	ta.startMtx.Lock()
	ta.windowActive = false
	ta.startMtx.Unlock()
}

func NewTimeBlockMaximum(windowSize int) Aggregator {
	return &TimeBlockMaximum{
		windowSize: time.Duration(windowSize),
		windowMtx:  &sync.Mutex{},
		startMtx:   &sync.RWMutex{},
	}
}
