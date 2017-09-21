package splitString

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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
	tc.SetInput("stringToSplit", "device/DEV01")
	tc.SetInput("delimiter", "/")
	tc.SetInput("tokenNumber", 2)

	act.Eval(tc)

	//check result attr
	value, _ := tc.GetOutput("token").(string)
	assert.Equal(t, "DEV01", value, "Bad Token value - TEST 1")

	//setup attrs
	tc.SetInput("stringToSplit", "device/DEV01")
	tc.SetInput("delimiter", "/")
	tc.SetInput("tokenNumber", 3)

	act.Eval(tc)

	//check result attr
	value, _ = tc.GetOutput("token").(string)
	assert.Equal(t, "", value, "Bad Token value - TEST 2")

	//setup attrs
	tc.SetInput("stringToSplit", "device")
	tc.SetInput("delimiter", "/")
	tc.SetInput("tokenNumber", 1)

	act.Eval(tc)

	//check result attr
	value, _ = tc.GetOutput("token").(string)
	assert.Equal(t, "device", value, "Bad Token value - TEST 3")

	//setup attrs
	tc.SetInput("stringToSplit", "")
	tc.SetInput("delimiter", "/")
	tc.SetInput("tokenNumber", 1)

	act.Eval(tc)

	//check result attr
	value, _ = tc.GetOutput("token").(string)
	assert.Equal(t, "", value, "Bad Token value - TEST 4")

	//setup attrs
	tc.SetInput("stringToSplit", "device/DEV01")
	tc.SetInput("delimiter", "/")
	tc.SetInput("tokenNumber", 0)

	act.Eval(tc)

	//check result attr
	value, _ = tc.GetOutput("token").(string)
	assert.Equal(t, "", value, "Bad Token value - TEST 5")

}
