package datatojson

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("activity-ecarlier-datatojson")

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
	obj := context.GetInput("input").(map[string]interface{})

	jsonBytes, err := json.Marshal(obj)

	if err != nil {
		log.Errorf("JSON Marshal error on object [%v]", obj)
		return false, err
	}

	context.SetOutput("result", string(jsonBytes))

	return true, nil
}
