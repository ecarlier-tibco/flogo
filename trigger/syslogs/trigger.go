package syslogs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	syslog "gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

// log is the default package logger
var log = logger.GetLogger("trigger-syslogs")

// MyTriggerFactory My Syslogs Trigger factory
type MyTriggerFactory struct {
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MyTriggerFactory{metadata: md}
}

// New -> Creates a new trigger instance for a given id
func (t *MyTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &MyTrigger{metadata: t.metadata, config: config}
}

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	server   *syslog.Server
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {
	// start the trigger

	if t.config.Settings == nil {
		panic(fmt.Sprintf("No Settings found for trigger '%s'", t.config.Id))
	}

	port := "514"

	if _, ok := t.config.Settings["port"]; !ok {
		log.Infof("No Port found for syslogs trigger '%s' in settings. Using default 514", t.config.Id)
	} else {
		port = t.config.GetSetting("port")
	}

	udp := true
	if _, ok := t.config.Settings["UDP"]; ok {
		udpConfig, err := strconv.ParseBool(t.config.GetSetting("UDP"))
		if err != nil {
			log.Warnf("Value set for UDP '%s' is not a boolean. Using default 'true'. Trigger '%s'",
				t.config.GetSetting("UDP"),
				t.config.Id)
		} else {
			udp = udpConfig
		}
	}
	tcp := false
	if _, ok := t.config.Settings["TCP"]; ok {
		tcpConfig, err := strconv.ParseBool(t.config.GetSetting("TCP"))
		if err != nil {
			log.Warnf("Value set for TCP '%s' is not a boolean. Using default 'false'. Trigger '%s'",
				t.config.GetSetting("TCP"),
				t.config.Id)
		} else {
			tcp = tcpConfig
		}
	}

	if !udp && !tcp {
		panic(fmt.Sprintf("At least one of UDP or TCP needs to be enabled for trigger '%s'", t.config.Id))
	}

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)

	for _, handlerCfg := range t.config.Handlers {
		handler := new(MySyslogHandler)
		handler.t = t
		handler.ActionID = handlerCfg.ActionId
		handler.handlerCfg = handlerCfg
		server.SetHandler(handler)
	}

	if udp {
		server.ListenUDP("0.0.0.0:" + port)
	}
	if tcp {
		server.ListenTCP("0.0.0.0:" + port)
	}

	server.Boot()

	t.server = server
	server.Wait()
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	t.server.Kill()
	return nil
}

// MySyslogHandler : My trigger's syslogs handler
type MySyslogHandler struct {
	t          *MyTrigger
	ActionID   string
	handlerCfg *trigger.HandlerConfig
}

// Handle : processing of each received syslogs
func (s *MySyslogHandler) Handle(logParts format.LogParts, msgLen int64, err error) {

	if err != nil {
		log.Errorf("Syslogs Handler for ActionID [%s] received error [%s]", s.ActionID, err)
	}

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	for key, value := range logParts {
		if key != "timestamp" {
			data[key] = value
		} else {
			ts, ok := value.(time.Time)
			if ok {
				data[key] = ts.Format(time.RFC3339)
			}
		}
	}
	req.Data = data

	md := s.t.Metadata()
	startAttrs, errorAttrs := md.OutputsToAttrs(req.Data, false)
	/*
		for _, attr := range startAttrs {
			spew.Dump(&attr)
		}
	*/
	if errorAttrs != nil || startAttrs == nil {
		log.Errorf("Failed to create output attributes for syslogs message for ActionID [%s] for reason [%s] message lost", s.ActionID, errorAttrs)
	}

	ctx := trigger.NewInitialContext(startAttrs, s.handlerCfg)

	action := action.Get(s.ActionID)
	_, _, errRun := s.t.runner.Run(ctx, action, s.ActionID, nil)

	if errRun != nil {
		log.Errorf("Failed to process syslogs message for ActionID [%s] for reason [%s] message lost", s.ActionID, errRun)
	}
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI string                 `json:"flowUri"`
	Data       map[string]interface{} `json:"data"`
	ReplyTo    string                 `json:"replyTo"`
}
