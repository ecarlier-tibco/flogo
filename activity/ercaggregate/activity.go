package ercaggregate

import (
	"errors"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/ecarlier-tibco/flogo/activity/ercaggregate/aggregator"
)

// activityLogger is the default logger for the Aggregate Activity
var activityLogger = logger.GetLogger("activity-tibco-aggregate")

const (
	ivFunction   = "functions"
	ivWindowSize = "windowSize"
	ivValue      = "value"

	ovResult = "result"
	ovReport = "report"
)

func init() {
	activityLogger.SetLogLevel(logger.InfoLevel)
}

// AggregateActivity is an Activity that is used to Aggregate a message to the console
// inputs : {function, windowSize, autoRest, value}
// outputs: {result, report}
type AggregateActivity struct {
	metadata *activity.Metadata
	mutex    *sync.RWMutex

	// aggregators stateful map of aggregators
	aggregators map[string]aggregator.Aggregator
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AggregateActivity{metadata: metadata, aggregators: make(map[string]aggregator.Aggregator), mutex: &sync.RWMutex{}}
}

// Metadata returns the activity's metadata
func (a *AggregateActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Aggregates the Message
func (a *AggregateActivity) Eval(context activity.Context) (done bool, err error) {

	aggregatorKey := context.ActivityHost().Name() + ":" + context.Name()
	aggrNamesList, _ := context.GetInput(ivFunction).([]string)
	resultsList := make(map[string]float64)
	reportsList := make(map[string]bool)

	for _, aggrName := range aggrNamesList {

		a.mutex.RLock()
		//get aggregator for activity, assumes flow & task names are unique
		aggr, ok := a.aggregators[aggregatorKey+":"+aggrName]

		a.mutex.RUnlock()

		//if window not create for this flow, create it

		if !ok {

			//go doesn't have lock upgrades or try, so do same check again

			a.mutex.Lock()
			aggr, ok = a.aggregators[aggregatorKey+":"+aggrName]

			if !ok {
				windowSize, _ := context.GetInput(ivWindowSize).(int)

				factory := aggregator.GetFactory(aggrName)

				if factory == nil {
					return false, errors.New("Unknown aggregator: " + aggrName)
				}

				aggr = factory(windowSize)
				a.aggregators[aggregatorKey+":"+aggrName] = aggr

				activityLogger.Debug("Aggregator created for ", aggregatorKey)
			}

			a.mutex.Unlock()
		}

		value, ok := context.GetInput(ivValue).(float64)

		if !ok {
			value, _ = data.CoerceToNumber(context.GetInput(ivValue))
		}

		reportsList[aggrName], resultsList[aggrName] = aggr.Add(value)
		activityLogger.Infof("RESULTS :\n!%v!\n!%v!", reportsList, resultsList)
	}

	context.SetOutput(ovReport, reportsList)
	context.SetOutput(ovResult, resultsList)
	activityLogger.Infof("END RESULTS :\n!%v!\n!%v!", reportsList, resultsList)

	return true, nil
}
