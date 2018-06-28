package aggregator

import (
	"sync"
	"time"
)

type TimeBlockMinimum struct {
	windowSize   time.Duration
	values       []float64
	windowMtx    *sync.Mutex
	startMtx     *sync.RWMutex
	windowActive bool
}

func init() {
	RegisterFactory("timeblockmin", NewTimeBlockMinimum)
}

func (ta *TimeBlockMinimum) Add(value float64) (bool, float64) {

	ta.windowMtx.Lock()
	ta.values = append(ta.values, value)
	ta.windowMtx.Unlock()

	if ta.startWindow() {
		time.Sleep(ta.windowSize * time.Millisecond)
		return true, ta.min()
	} else {
		return false, 0
	}
}

func (ta *TimeBlockMinimum) min() float64 {

	ta.windowMtx.Lock()

	var min = float64(0)

	count := len(ta.values)

	for i := 0; i < count; i++ {
		if min == 0 || ta.values[i] < min {
			min = ta.values[i]
		}
	}

	ta.resetWindow()

	return min
}

func (ta *TimeBlockMinimum) startWindow() bool {

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

func (ta *TimeBlockMinimum) resetWindow() {
	ta.values = nil
	ta.windowMtx.Unlock()

	ta.startMtx.Lock()
	ta.windowActive = false
	ta.startMtx.Unlock()
}

func NewTimeBlockMinimum(windowSize int) Aggregator {
	return &TimeBlockMinimum{
		windowSize: time.Duration(windowSize),
		windowMtx:  &sync.Mutex{},
		startMtx:   &sync.RWMutex{},
	}
}
