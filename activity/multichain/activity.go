package multichain

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	multichain "github.com/golangdaddy/multichain-client"
)

var log = logger.GetLogger("activity-multichain")

type jsonPublishData struct {
	JSON interface{} `json:"json"`
}

const (
	ivChain       = "chain"
	ivHost        = "host"
	ivPort        = "port"
	ivRPCUser     = "rpcuser"
	ivRPCPassword = "rpcpassword"
	ivCommand     = "command"
	ivParameters  = "parameters"
	ivJSONData    = "jsonData"
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

	loglevel := log.GetLogLevel()
	if loglevel == logger.DebugLevel {
		a.client = a.client.DebugMode()
	}

	cmd := context.GetInput(ivCommand).(string)

	switch cmd {
	case "getinfo":
		log.Debug("Sending getinfo command")
		resp, err := a.client.GetInfo()
		if err != nil {
			log.Errorf("multichain getinfo error [%v]", err)
			return false, err
		}
		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("getinfo response [%v]", resp["result"])

	case "getaddresses":
		log.Debug("Sending getaddresses command")
		resp, err := a.client.GetAddresses(false)
		if err != nil {
			log.Errorf("multichain getaddresses error [%v]", err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("getaddresses response [%v]", resp["result"])

	case "create":
		// type, name, open (bool) parameters
		log.Debug("Sending create command")

		params := context.GetInput(ivParameters).(map[string]interface{})

		var ok, open bool
		var typeToCreate, name, openstr string

		typeToCreate, ok = params["type"].(string)
		if !ok {
			log.Error("Missing mandatory parameter [type] for [create] command")
			return false, nil
		}

		name, ok = params["name"].(string)
		if !ok {
			log.Error("Missing mandatory parameter [name] for [create] command")
			return false, nil
		}

		open = false
		openstr, ok = params["open"].(string)
		if ok && openstr == "true" {
			open = true
		}

		msg := a.client.Command(
			"create",
			[]interface{}{
				typeToCreate,
				name,
				open,
			},
		)

		resp, err := a.client.Post(msg)

		if err != nil {
			log.Errorf("multichain create error [%v]", err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("create response [%v]", resp["result"])

	case "grant":
		// addresses, permissions list of strings (comma-separated) parameters
		log.Debug("Sending grant command")

		params := context.GetInput(ivParameters).(map[string]interface{})

		var ok bool
		var addressesStr, permissionsStr string

		addressesStr, ok = params["addresses"].(string)
		if !ok {
			log.Error("Missing mandatory parameter [addresses] for [grant] command")
			return false, nil
		}

		addresses := strings.Split(addressesStr, ",")

		permissionsStr, ok = params["permissions"].(string)
		if !ok {
			log.Error("Missing mandatory parameter [permissions] for [grant] command")
			return false, nil
		}

		permissions := strings.Split(permissionsStr, ",")

		resp, err := a.client.Grant(addresses, permissions)

		if err != nil {
			log.Errorf("multichain grant error [%v]", err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("grant response [%v]", resp["result"])

	case "publish":
		log.Debug("Sending publish command")
		params := context.GetInput(ivParameters).(map[string]interface{})

		var ok bool
		var stream, key, data string

		stream, ok = params["stream"].(string)
		if !ok {
			log.Error("Missing mandatory parameter [stream] for [publish] command")
			return false, nil
		}

		log.Debugf("Stream Name [%v]", params["stream"])

		key, ok = params["key"].(string)
		if !ok {
			log.Error("Missing mandatory parameter [key] for [publish] command")
			return false, nil
		}

		data, ok = params["data"].(string)
		var jsonData interface{}
		if !ok {
			jsonData = context.GetInput(ivJSONData)

			if jsonData == nil {
				log.Error("Missing mandatory parameter [data] for [publish] command")
				return false, nil
			}
			log.Debugf("JSON DATA [%v]", jsonData)
		} else {
			log.Debugf("%s", data)
		}
		from, fromPresent := params["from"].(string)
		if fromPresent {
			log.Debugf("Received From [%v]", from)

			var msg interface{}
			if jsonData != nil {
				dataobj := &jsonPublishData{
					JSON: jsonData,
				}
				msg = a.client.Command(
					"publishfrom",
					[]interface{}{
						from,
						stream,
						key,
						dataobj,
					},
				)
			} else {
				msg = a.client.Command(
					"publishfrom",
					[]interface{}{
						from,
						stream,
						key,
						data,
					},
				)
			}

			resp, err := a.client.Post(msg)

			if err != nil {
				log.Errorf("multichain publishfrom error [%v]", err)
				return false, err
			}

			context.SetOutput(ovSuccess, true)
			context.SetOutput(ovResponse, resp["result"])
			log.Debugf("publishfrom response [%v]", resp["result"])
		} else {
			log.Debug("No from defined")

			var msg interface{}
			if jsonData != nil {
				dataobj := &jsonPublishData{
					JSON: jsonData,
				}
				msg = a.client.Command(
					"publish",
					[]interface{}{
						stream,
						key,
						dataobj,
					},
				)
			} else {
				msg = a.client.Command(
					"publish",
					[]interface{}{
						stream,
						key,
						data,
					},
				)
			}

			resp, err := a.client.Post(msg)

			if err != nil {
				log.Errorf("multichain publish error [%v]", err)
				return false, err
			}

			context.SetOutput(ovSuccess, true)
			context.SetOutput(ovResponse, resp["result"])
			log.Debugf("publish response [%v]", resp["result"])

		}

	case "subscribe", "unsubscibe":
		log.Debugf("Sending %v command", cmd)
		params := context.GetInput(ivParameters).(map[string]interface{})

		var ok bool
		var stream string

		stream, ok = params["stream"].(string)
		if !ok {
			log.Errorf("Missing mandatory parameter [stream] for [%v] command", cmd)
			return false, nil
		}

		msg := a.client.Command(
			cmd,
			[]interface{}{
				stream,
			},
		)
		resp, err := a.client.Post(msg)

		if err != nil {
			log.Errorf("multichain %v error [%v]", cmd, err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("%v response [%v]", cmd, resp["result"])

	case "getstreamitem":
		log.Debugf("Sending %v command", cmd)
		params := context.GetInput(ivParameters).(map[string]interface{})

		var ok bool
		var stream, txid string

		stream, ok = params["stream"].(string)
		if !ok {
			log.Errorf("Missing mandatory parameter [stream] for [%v] command", cmd)
			return false, nil
		}

		txid, ok = params["txid"].(string)
		if !ok {
			log.Errorf("Missing mandatory parameter [txid] for [%v] command", cmd)
			return false, nil
		}

		msg := a.client.Command(
			cmd,
			[]interface{}{
				stream,
				txid,
			},
		)
		resp, err := a.client.Post(msg)

		if err != nil {
			log.Errorf("multichain %v error [%v]", cmd, err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("%v response [%v]", cmd, resp["result"])

	case "liststreamkeyitems":
		log.Debugf("Sending %v command", cmd)
		params := context.GetInput(ivParameters).(map[string]interface{})

		var ok bool
		var stream, key string

		stream, ok = params["stream"].(string)
		if !ok {
			log.Errorf("Missing mandatory parameter [stream] for [%v] command", cmd)
			return false, nil
		}

		key, ok = params["key"].(string)
		if !ok {
			log.Errorf("Missing mandatory parameter [key] for [%v] command", cmd)
			return false, nil
		}

		msg := a.client.Command(
			cmd,
			[]interface{}{
				stream,
				key,
			},
		)
		resp, err := a.client.Post(msg)

		if err != nil {
			log.Errorf("multichain %v error [%v]", cmd, err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("%v response [%v]", cmd, resp["result"])

	case "liststreamkeys", "liststreamitems", "liststreampublishers":
		log.Debugf("Sending %v command", cmd)
		params := context.GetInput(ivParameters).(map[string]interface{})

		var ok bool
		var stream string

		stream, ok = params["stream"].(string)
		if !ok {
			log.Errorf("Missing mandatory parameter [stream] for [%v] command", cmd)
			return false, nil
		}

		msg := a.client.Command(
			cmd,
			[]interface{}{
				stream,
			},
		)
		resp, err := a.client.Post(msg)

		if err != nil {
			log.Errorf("multichain %v error [%v]", cmd, err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("%v response [%v]", cmd, resp["result"])

	case "liststreampublisheritems":
		log.Debugf("Sending %v command", cmd)
		params := context.GetInput(ivParameters).(map[string]interface{})

		var ok bool
		var stream, address string

		stream, ok = params["stream"].(string)
		if !ok {
			log.Errorf("Missing mandatory parameter [stream] for [%v] command", cmd)
			return false, nil
		}

		address, ok = params["address"].(string)
		if !ok {
			log.Errorf("Missing mandatory parameter [address] for [%v] command", cmd)
			return false, nil
		}

		msg := a.client.Command(
			cmd,
			[]interface{}{
				stream,
				address,
			},
		)
		resp, err := a.client.Post(msg)

		if err != nil {
			log.Errorf("multichain %v error [%v]", cmd, err)
			return false, err
		}

		context.SetOutput(ovSuccess, true)
		context.SetOutput(ovResponse, resp["result"])
		log.Debugf("%v response [%v]", cmd, resp["result"])

	default:
		log.Errorf("Invalid Command received [%v]", cmd)
		return false, nil

	}

	return true, nil
}
