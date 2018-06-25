package multichain

import (
	"encoding/json"
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
	tc.SetInput("rpcuser", "yourrpcuser")
	tc.SetInput("rpcpassword", "yourrpcpassword")

	var ok bool

	tc.SetInput("command", "getinfo")
	act.Eval(tc)

	ok = tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("getinfo response [%v]", tc.GetOutput("response"))
	}

	tc.SetInput("command", "getaddresses")
	act.Eval(tc)

	ok = tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("getaddresses response [%v]", tc.GetOutput("response"))
	}

	var parameters map[string]string

	tc.SetInput("command", "publish")

	parameters = make(map[string]string)
	parameters["stream"] = "teststream"
	parameters["key"] = "testkey1"
	tc.SetInput("parameters", parameters)

	var dataobj map[string]interface{}
	json.Unmarshal([]byte(`{"i":[98,99],"j":"sure"}`), &dataobj)

	tc.SetInput("jsonData", dataobj)

	act.Eval(tc)

	var txid string
	ok = tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("publish response [%v]\n", tc.GetOutput("response"))
		txid = tc.GetOutput("response").(string)
	}

	tc.SetInput("command", "getstreamitem")
	parameters = make(map[string]string)
	parameters["stream"] = "teststream"
	parameters["txid"] = txid

	tc.SetInput("parameters", parameters)

	act.Eval(tc)

	ok = tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("getstreamitem response [%v]\n", tc.GetOutput("response"))
	}

	// liststreamkeys test
	tc.SetInput("command", "liststreamkeys")
	parameters = make(map[string]string)
	parameters["stream"] = "teststream"

	tc.SetInput("parameters", parameters)

	act.Eval(tc)

	ok = tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("liststreamkeys response [%v]\n", tc.GetOutput("response"))
	}

	// test liststreamkeyitems
	tc.SetInput("command", "liststreamkeyitems")
	parameters = make(map[string]string)
	parameters["stream"] = "teststream"
	parameters["key"] = "testkey1"
	tc.SetInput("parameters", parameters)

	act.Eval(tc)

	ok = tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("liststreamkeyitems response [%v]\n", tc.GetOutput("response"))
	}

	// test liststreamkeyitems
	tc.SetInput("command", "liststreampublishers")
	parameters = make(map[string]string)
	parameters["stream"] = "teststream"
	tc.SetInput("parameters", parameters)

	act.Eval(tc)

	ok = tc.GetOutput("success").(bool)
	var publisher string
	if ok {
		fmt.Printf("liststreampublishers response [%v]\n", tc.GetOutput("response"))
		resparray := tc.GetOutput("response").([]interface{})
		respitem := resparray[0].(map[string]interface{})
		publisher = respitem["publisher"].(string)

	}

	// 17uyvHrXm3XRp68eJTtwi2pBTVVrHrSfmNVCrn
	// test liststreampublisheritems
	tc.SetInput("command", "liststreampublisheritems")
	parameters = make(map[string]string)
	parameters["stream"] = "teststream"
	parameters["address"] = publisher
	tc.SetInput("parameters", parameters)

	act.Eval(tc)

	ok = tc.GetOutput("success").(bool)
	if ok {
		fmt.Printf("liststreampublisheritems response [%v]\n", tc.GetOutput("response"))
	}

}
