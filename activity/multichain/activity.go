package multichain

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	multichain "github.com/golangdaddy/multichain-client"
)

var log = logger.GetLogger("activity-multichain")

const (
	ivChain       = "chain"
	ivHost        = "host"
	ivPort        = "port"
	ivRPCUser     = "rpcuser"
	ivRPCPassword = "rpcpassword"
	ivCommand     = "command"
	ovSuccess     = "success"
	ovResponse    = "response"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
	client   *multichain.Client
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

	// do eval

	if a.client == nil {
		log.Info("Connecting to MultiChain")
		host := context.GetInput(ivHost).(string)
		port := context.GetInput(ivPort).(int)
		rpcuser := context.GetInput(ivRPCUser).(string)
		rpcpassword := context.GetInput(ivRPCPassword).(string)
		chain := context.GetInput(ivChain).(string)
		a.client = multichain.NewClient(chain, rpcuser, rpcpassword, port)
		a.client = a.client.ViaNode(host, port)
	}
	a.client = a.client.DebugMode()

	cmd := context.GetInput(ivCommand).(string)
	switch cmd {
	case "getinfo":
		resp, err := a.client.GetInfo()
		if err != nil {
			log.Errorf("multichain getinfo error [%v]", err)
			return false, err
		}
		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])

	case "getaddresses":
		resp, err := a.client.GetAddresses(false)
		if err != nil {
			log.Errorf("multichain getaddresses error [%v]", err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])

	default:
		log.Errorf("Invalid Command received [%v]", cmd)
		return false, nil

	}

	return true, nil
}
