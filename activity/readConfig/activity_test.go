package readConfig

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivFileName, "./testConfig.json")

	act.Eval(tc)

	//check result attr
	value, _ := tc.GetOutput(ovResult).(string)
	assert.Equal(t, "ok", value, "Bad Result")

	typedVal, ok := data.GetGlobalScope().GetAttr("ParamStr")
	assert.True(t, ok, "Get Global Scope on string parameter Failed")
	assert.Equal(t, "ValueStr", typedVal.Value, "Wrong value for string parameter")

	typedVal, ok = data.GetGlobalScope().GetAttr("ParamInt")
	assert.True(t, ok, "Get Global Scope on integer parameter Failed")
	assert.Equal(t, 22.0, typedVal.Value, "Wrong value for integer parameter")

	typedVal, ok = data.GetGlobalScope().GetAttr("ParamNumber")
	assert.True(t, ok, "Get Global Scope on number parameter Failed")
	assert.Equal(t, 23.897, typedVal.Value.(float64), "Wrong value for number parameter")

	typedVal, ok = data.GetGlobalScope().GetAttr("ParamBool")
	assert.True(t, ok, "Get Global Scope on number parameter Failed")
	assert.True(t, typedVal.Value.(bool), "Wrong value for boolean parameter")

}
