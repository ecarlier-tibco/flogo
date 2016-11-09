package pushover

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/op/go-logging"
	"github.com/toorop/pushover"
	"time"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-pushover")

const (
	ivUserKey = "userKey"
	ivAuthToken = "authToken"
	ivMessage = "message"
	ivMessageTitle = "messageTitle"
	ivDevice = "device"
	ivSound = "sound"
	ivURL = "url"
	ivURLTitle = "urlTitle"
	ivPriority = "priority"
	ovStatus = "status"
)

// MyActivity is a stub for your Activity implementation
type PushoverActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&PushoverActivity{metadata: md})
}

// Metadata implements activity.Activity.Metadata
func (a *PushoverActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *PushoverActivity) Eval(context activity.Context) (done bool, err error) {

	userKey := context.GetInput(ivUserKey)
	authToken := context.GetInput(ivAuthToken)
	message := context.GetInput(ivMessage)
	messageTitle := context.GetInput(ivMessageTitle)
	device := context.GetInput(ivDevice)
	sound := context.GetInput(ivSound)
	url := context.GetInput(ivURL)
	urlTitle := context.GetInput(ivURLTitle)
	priority := context.GetInput(ivPriority)

	// Check if mandatory credentials are set in config
	if userKey == nil || authToken == nil {
		log.Error("Missing Pushover Credentials")
		err := activity.NewError("Credential Config not specified")
		return false, err
	}

	// Check if there is a message to send
	if message == nil {
		log.Error("No Message to Send")
		context.SetOutput(ovStatus, "NO_MSG")
		return true, nil
	}

	// Now let's try to connect
	po, err := pushover.NewPushover(authToken.(string), userKey.(string))
	if err != nil {
		log.Error("PushOver Connection Error : ", err)
		context.SetOutput(ovStatus, "CONNECT_ERR")
		return true, nil
	}

	m := pushover.Message{}

	m.Message = message.(string)
	if priority != nil {
		m.Priority = priority.(int)
	}
	if messageTitle != nil {
		m.Title = messageTitle.(string)
	}
	if device != nil {
		m.Device = device.(string)
	}
	if url != nil {
		m.Url = url.(string)
	}
	if urlTitle != nil {
		m.UrlTitle = urlTitle.(string)
	}
	m.Timestamp = time.Now().Unix()
	if sound != nil {
		m.Sound = sound.(string)
	}

	if preq, presp, err := po.Push(&m); err != nil {
		log.Error("Push Message Error : ", err)
		log.Error("Sent : ", preq)
		log.Error("Received : ", presp)
		context.SetOutput(ovStatus, "PUSH_ERR")
		return true, nil
	}

	context.SetOutput(ovStatus, "OK")

	return true, nil
}
