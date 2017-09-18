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
	typedVal, ok := data.GetGlobalScope().GetAttr("Param1")
	assert.True(t, ok, "Get Global Scope Failed")
	assert.Equal(t, typedVal.Value, "Value1", "Wrong value")
	typedVal, ok = data.GetGlobalScope().GetAttr("Param2")
	assert.True(t, ok, "Get Global Scope Failed")
	assert.Equal(t, typedVal.Value, "Value2", "Wrong value")

}
