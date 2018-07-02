package stringtonumber

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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

	var in map[string]interface{}
	in = make(map[string]interface{})
	in["ID"] = "Test66"
	in["Chain_Address"] = "mc_address_6666"
	in["Latitude"] = 42.76
	in["Longitude"] = 4.55
	in["MeasureTypes"] = [2]string{"temperature", "humidity"}

	tc.SetInput("input", in)
	act.Eval(tc)
	result := tc.GetOutput("result").(string)
	fmt.Printf("GOT RESULT [%s]\n", result)

}
