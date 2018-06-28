package aggregator

import "sync"

type SizeBlock struct {
	windowSize   int
	values       []float64
	nextValueIdx int
	mutex        *sync.Mutex
	functions    []string
}

func init() {
	RegisterFactory("sizeblock", NewSizeBlock)
}

func (ba *SizeBlock) Add(value float64) (bool, map[string]interface{}) {

	ba.mutex.Lock()
	defer ba.mutex.Unlock()

	ba.values[ba.nextValueIdx] = value
	ba.nextValueIdx = ba.nextValueIdx + 1

	if ba.nextValueIdx >= ba.windowSize {
		return true, ba.calculateFunctions()
	}

	return false, nil
}

func (ba *SizeBlock) calculateFunctions() map[string]interface{} {

	var total = float64(0)
	var min = float64(0)
	var max = float64(0)

	for i := 0; i < ba.windowSize; i++ {
		total += ba.values[i]
		if max == 0 || ba.values[i] > max {
			max = ba.values[i]
		}
		if min == 0 || ba.values[i] < min {
			min = ba.values[i]
		}
	}
	ba.nextValueIdx = 0

	var results map[string]interface{}
	results = make(map[string]interface{})

	for _, f := range ba.functions {
		switch f {
		case "avg", "average":
			results[f] = total / float64(ba.windowSize)
		case "min", "minimum":
			results[f] = min
		case "max", "maximum":
			results[f] = max
		}
	}

	return results
}

func NewSizeBlock(windowSize int, functions []string) Aggregator {
	return &SizeBlock{
		windowSize: windowSize,
		values:     make([]float64, windowSize),
		mutex:      &sync.Mutex{},
		functions:  functions,
	}
}
