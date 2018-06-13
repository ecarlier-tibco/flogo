package numdatasim

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/agoussia/godes"
)

var log = logger.GetLogger("activity-num-data-generator")

const (
	ivDistributionType = "distributionType"
	ivMin              = "min"
	ivMax              = "max"
	ivMean             = "mean"
	ivSigma            = "sigma"
	ivMode             = "mode"
	ivLambda           = "lambda"
	ivRepeatSequence   = "repeatSequence"

	ovOutput = "output"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
	distrib  interface{}
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	distribType := context.GetInput(ivDistributionType).(string)
	repeatFlag := context.GetInput(ivRepeatSequence).(bool)

	switch distribType {
	case "Normal":
		mean := context.GetInput(ivMean).(float64)
		sigma := context.GetInput(ivSigma).(float64)
		log.Debugf("Normal Distribution with Mean [%v] and sigma [%v]", mean, sigma)

		if a.distrib == nil {
			log.Debugf("Creating Normal Distributor / Repeat = %v", repeatFlag)
			a.distrib = godes.NewNormalDistr(repeatFlag)
		}
		ldistrib := a.distrib.(*godes.NormalDistr)
		newdata := ldistrib.Get(mean, sigma)

		context.SetOutput(ovOutput, newdata)

	case "Uniform":
		min := context.GetInput(ivMin).(float64)
		max := context.GetInput(ivMax).(float64)
		log.Debugf("Uniform Distribution with Min [%v] and Max [%v]", min, max)

		if a.distrib == nil {
			log.Debugf("Creating Uniform Distributor / Repeat = %v", repeatFlag)
			a.distrib = godes.NewUniformDistr(repeatFlag)
		}
		ldistrib := a.distrib.(*godes.UniformDistr)
		newdata := ldistrib.Get(min, max)

		context.SetOutput(ovOutput, newdata)

	case "Triangular":
		min := context.GetInput(ivMin).(float64)
		max := context.GetInput(ivMax).(float64)
		mode := context.GetInput(ivMode).(float64)

		log.Debugf("Triangular Distribution with Min [%v], Mode [%v] and Max [%v]", min, mode, max)

		if a.distrib == nil {
			log.Debugf("Creating Triangular Distributor / Repeat = %v", repeatFlag)
			a.distrib = godes.NewTriangularDistr(repeatFlag)
		}
		ldistrib := a.distrib.(*godes.TriangularDistr)
		newdata := ldistrib.Get(min, max, mode)

		context.SetOutput(ovOutput, newdata)

	case "Exponential":
		lambda := context.GetInput(ivLambda).(float64)

		log.Debugf("Exponential Distribution with Lamda [%v]", lambda)

		if a.distrib == nil {
			log.Debugf("Creating Exponential Distributor / Repeat = %v", repeatFlag)
			a.distrib = godes.NewExpDistr(repeatFlag)
		}
		ldistrib := a.distrib.(*godes.ExpDistr)
		newdata := ldistrib.Get(lambda)

		context.SetOutput(ovOutput, newdata)

	default:
		log.Errorf("Invalid Distribution Type [%v]", distribType)
		return false, nil
	}
	return true, nil
}
