package aggregator

import (
	"sync"
	"time"
)

type TimeBlock struct {
	windowSize   time.Duration
	values       []float64
	windowMtx    *sync.Mutex
	startMtx     *sync.RWMutex
	windowActive bool
	functions    []string
}

func init() {
	RegisterFactory("timeblock", NewTimeBlock)
}

func (ta *TimeBlock) Add(value float64) (bool, map[string]interface{}) {

	ta.windowMtx.Lock()
	ta.values = append(ta.values, value)
	ta.windowMtx.Unlock()

	if ta.startWindow() {
		time.Sleep(ta.windowSize * time.Millisecond)
		return true, ta.calculateFunctions()
	} else {
		return false, nil
	}
}

func (ta *TimeBlock) calculateFunctions() map[string]interface{} {

	ta.windowMtx.Lock()

	var total = float64(0)
	var min = float64(0)
	var max = float64(0)

	count := len(ta.values)

	for i := 0; i < count; i++ {
		total += ta.values[i]
		if max == 0 || ta.values[i] > max {
			max = ta.values[i]
		}
		if min == 0 || ta.values[i] < min {
			min = ta.values[i]
		}
	}

	ta.resetWindow()

	var results map[string]interface{}
	results = make(map[string]interface{})

	for _, f := range ta.functions {
		switch f {
		case "avg", "average":
			results[f] = total / float64(count)
		case "min", "minimum":
			results[f] = min
		case "max", "maximum":
			results[f] = max
		}
	}

	return results
}

func (ta *TimeBlock) startWindow() bool {

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

func (ta *TimeBlock) resetWindow() {
	ta.values = nil
	ta.windowMtx.Unlock()

	ta.startMtx.Lock()
	ta.windowActive = false
	ta.startMtx.Unlock()
}

func NewTimeBlock(windowSize int, functions []string) Aggregator {
	return &TimeBlock{
		windowSize: time.Duration(windowSize),
		windowMtx:  &sync.Mutex{},
		startMtx:   &sync.RWMutex{},
		functions:  functions,
	}
}
