package aggregator

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

type Aggregator interface {
	Add(value float64) (report bool, result map[string]interface{})
}

type Factory func(windowSize int, functions []string) Aggregator

var (
	factoryMu sync.RWMutex
	factories = make(map[string]Factory)
)

// Register registers the specified aggregator
func RegisterFactory(id string, factory Factory) {
	factoryMu.Lock()
	defer factoryMu.Unlock()

	logger.Debugf("Registering factory: '%s'", id)

	if factory == nil {
		panic("aggregator.RegisterFactory: factory is nil")
	}

	if _, dup := factories[id]; dup {
		panic("aggregator.Register: aggregator already registered " + id)
	}

	factories[id] = factory
}

// Get gets specified aggregator factory
func GetFactory(id string) Factory {
	return factories[id]
}
