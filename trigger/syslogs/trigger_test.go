package syslogs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/davecgh/go-spew/spew"
)

// var exampleRFC5424Syslog = "<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47 - 'su root' failed for lonvick on /dev/pts/8"
// var extremeSyslogs1 = "<142>events: EventType[Roam] MAC[78:D7:5F:80:B8:87] AP[AP-Reunion] FromAP[AP-Reunion] BSSID[20:B3:99:D8:F2:00] Details: Inside AC from AP/Radio[2] to AP/Radio[1] VNS[Tibco IoT]"
var extremeSyslogs2 = "<27>Jun 27 14:21:41 purview appid: LEEF:1.0|Extreme Networks|Purview|1.0|Purview Flow|devTime=1467030101	proto=TCP	src=192.168.171.17	dst=17.252.11.246	srcPort=50951	dstPort=443	Application=Apple	ApplicationInt=12975	ClientOS=iOS iPhone	ClientOSFamily=iOS	DHCP_ClientIP=192.168.171.17	ServerIP=17.252.11.246	SwitchType=CoreFlow	uuid=a0379a36	"

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	fmt.Println("TestRunner")
	fmt.Printf("URI: [%s]", uri)
	spew.Dump(context)

	return 0, nil, nil
}

const testConfig string = `{
  "name": "syslogs",	
  "settings": {
		"port": "11514",
		"UDP": "true",
		"TCP": "false"
  },
  "handlers": [
    {
      "actionId": "test_action"
    }
  ]
}`

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
func TestInit(t *testing.T) {

	// New factory
	md := trigger.NewMetadata(getJSONMetadata())
	f := NewFactory(md)

	// New Trigger
	config := trigger.Config{}

	fmt.Println(testConfig)

	json.Unmarshal([]byte(testConfig), &config)
	tgr := f.New(&config)

	runner := &TestRunner{}

	go func() {
		tgr.Init(runner)
		tgr.Start()
	}()

	time.Sleep(time.Second * 1)
	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:11514")
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()
	// Conn.Write([]byte(exampleRFC5424Syslog))
	Conn.Write([]byte(extremeSyslogs2))
	time.Sleep(time.Second * 1)

	tgr.Stop()

}
