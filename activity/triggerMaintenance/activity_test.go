/*
 * Copyright © 2017. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package triggerMaintenance

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

const (
	user = "ecarlier@tibco.com"
	pass = "Leo300307"
)

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

func TestActivityRegistration(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("username", "[YOUR_LIVEAPPS_USERNAME]")
	tc.SetInput("password", "[YOUR_LIVEAPPS_PASSWORD]")
	tc.SetInput("region", "EU")
	tc.SetInput("accountID", "01BJRKPJCNH0G63DC8YRXVNF3J")
	tc.SetInput("sandboxID", "342")
	tc.SetInput("applicationID", 511)
	tc.SetInput("applicationCreatorID", 2302)
	tc.SetInput("applicationVersion", 1)
	tc.SetInput("alert", "Alert from GOTEST")

	dt, _ := data.ToTypeEnum("string")
	data.GetGlobalScope().AddAttr("DeviceID", dt, "DEV_GOTEST")
	dt, _ = data.ToTypeEnum("number")
	data.GetGlobalScope().AddAttr("DeviceLatitude", dt, 48.875621)
	data.GetGlobalScope().AddAttr("DeviceLongitude", dt, 2.302543)

	_, err := act.Eval(tc)
	assert.Nil(t, err)
	result := tc.GetOutput("result")

	fmt.Println("Result is:" + result.(string))
	assert.NotNil(t, result)
	// use this to see output on a valid test
	// assert.NotNil(t, nil)
}
