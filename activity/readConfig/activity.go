package readConfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("activity-ecarlier-readConfig")

const (
	ivFileName = "configFile"

	ovResult = "result"
)

type jsonobject struct {
	Config ConfigType
}

type ConfigType struct {
	Parameters []ParameterType
}

type ParameterType struct {
	Name  string
	Type  string
	Value interface{}
}

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

	fileName := context.GetInput(ivFileName).(string)
	fileContent, err := ioutil.ReadFile(fileName)

	if err != nil {
		errorMsg := fmt.Sprintf("Could not open file for read '%s'", fileName)
		log.Error(errorMsg)
		return false, activity.NewError(errorMsg, "", nil)
	}

	var jsonconfig jsonobject
	err = json.Unmarshal(fileContent, &jsonconfig)
	if err != nil {
		errorMsg := fmt.Sprintf("JSON Unmarshal error '%s'", err.Error())
		log.Error(errorMsg)
		return false, activity.NewError(errorMsg, "", nil)
	}
	for _, param := range jsonconfig.Config.Parameters {
		dt, ok := data.ToTypeEnum(param.Type)
		if ok {
			data.GetGlobalScope().AddAttr(param.Name, dt, param.Value)
		} else {
			log.Errorf("Wrong type %s specified for Config Parameter %s", param.Type, param.Name)
		}
	}

	context.SetOutput(ovResult, "ok")
	return true, nil
}
