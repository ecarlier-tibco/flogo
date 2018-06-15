package multichain

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

	//setup attrs
	tc.SetInput("chain", "demochain")
	tc.SetInput("host", "localhost")
	tc.SetInput("port", 6284)
	tc.SetInput("rpcuser", "tibco")
	tc.SetInput("rpcpassword", "balzac")

	tc.SetInput("command", "getinfo")
	act.Eval(tc)

	//check result attr
	ok := tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("getinfo response [%v]", tc.GetOutput("response"))
	}

	tc.SetInput("command", "getaddresses")
	act.Eval(tc)

	//check result attr
	ok = tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("getaddresses response [%v]", tc.GetOutput("response"))
	}

}
