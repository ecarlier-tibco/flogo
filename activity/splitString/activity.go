package splitString

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	ivStringToSplit = "stringToSplit"
	ivDelimiter     = "delimiter"
	ivTokenNumber   = "tokenNumber"
	ovToken         = "token"
)

// log is the default package logger
var log = logger.GetLogger("activity-ecarlier-splitstring")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
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
	stringToSplit := context.GetInput(ivStringToSplit).(string)
	delimiter := context.GetInput(ivDelimiter).(string)
	tokenNumber := context.GetInput(ivTokenNumber).(int)

	if tokenNumber < 1 {
		log.Error("Bad token number value. Should be >=1")
		context.SetOutput(ovToken, "")
		return true, nil
	}

	tokens := strings.Split(stringToSplit, delimiter)
	numtokens := len(tokens)

	if tokenNumber > numtokens {
		log.Info("Token Number higher than number of tokens")
		context.SetOutput(ovToken, "")
		return true, nil
	}

	context.SetOutput(ovToken, tokens[tokenNumber-1])
	return true, nil
}
