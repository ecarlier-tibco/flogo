package aggregator

import "sync"

type MovingBlock struct {
	windowSize   int
	values       []float64
	nextValueIdx int
	full         bool
	mutex        *sync.Mutex
	functions    []string
}

func init() {
	RegisterFactory("movingblock", NewMovingBlock)
}

func (ma *MovingBlock) Add(value float64) (bool, map[string]interface{}) {

	ma.mutex.Lock()
	defer ma.mutex.Unlock()

	//if ma.full && ma.nextValueIdx == 0 && ma.autoReset {
	//	ma.full = false
	//}

	ma.values[ma.nextValueIdx] = value

	ma.nextValueIdx = (ma.nextValueIdx + 1) % ma.windowSize

	if !ma.full && ma.nextValueIdx == 0 {
		ma.full = true
	}

	if ma.full {
		return true, ma.calculateFunctions()
	}

	return false, nil
}

func (ma *MovingBlock) calculateFunctions() map[string]interface{} {

	var total = float64(0)
	var min = float64(0)
	var max = float64(0)

	var count = ma.windowSize

	if !ma.full {
		if ma.nextValueIdx == 0 {
			return nil
		}

		count = ma.nextValueIdx
	}

	for i := 0; i < count; i++ {
		total += ma.values[i]
		if max == 0 || ma.values[i] > max {
			max = ma.values[i]
		}
		if min == 0 || ma.values[i] < min {
			min = ma.values[i]
		}
	}

	var results map[string]interface{}
	results = make(map[string]interface{})

	for _, f := range ma.functions {
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

func NewMovingBlock(windowSize int, functions []string) Aggregator {
	return &MovingBlock{
		windowSize: windowSize,
		values:     make([]float64, windowSize),
		mutex:      &sync.Mutex{},
		functions:  functions,
	}
}
