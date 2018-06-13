package numdatasim

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

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

	// 1 Test Normal Distribution
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("distributionType", "Normal")
	tc.SetInput("mean", 10.0)
	tc.SetInput("sigma", 1.0)

	for i := 0; i < 5; i++ {
		act.Eval(tc)

		//check result attr
		val := tc.GetOutput("output")

		fmt.Printf("Normal Distribution output: %v\n", val)
		time.Sleep(1 * time.Second)
	}

	// 2 Test Uniform Distribution
	act = NewActivity(getActivityMetadata())
	tc = test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("distributionType", "Uniform")
	tc.SetInput("min", 10.0)
	tc.SetInput("max", 15.0)
	tc.SetInput("repeatSequence", true)

	for i := 0; i < 5; i++ {
		act.Eval(tc)

		//check result attr
		val := tc.GetOutput("output")

		fmt.Printf("Uniform Distribution output: %v\n", val)
		time.Sleep(1 * time.Second)
	}

	// 3 Test Triangular Distribution
	act = NewActivity(getActivityMetadata())
	tc = test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("distributionType", "Triangular")
	tc.SetInput("min", 10.0)
	tc.SetInput("max", 15.0)
	tc.SetInput("mode", 12.0)

	for i := 0; i < 5; i++ {
		act.Eval(tc)

		//check result attr
		val := tc.GetOutput("output")

		fmt.Printf("Triangular Distribution output: %v\n", val)
		time.Sleep(1 * time.Second)
	}

	// 4 Test Exponential Distribution
	act = NewActivity(getActivityMetadata())
	tc = test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("distributionType", "Exponential")
	tc.SetInput("lambda", 0.01)

	for i := 0; i < 5; i++ {
		act.Eval(tc)

		//check result attr
		val := tc.GetOutput("output")

		fmt.Printf("Triangular Distribution output: %v\n", val)
		time.Sleep(1 * time.Second)
	}
}
