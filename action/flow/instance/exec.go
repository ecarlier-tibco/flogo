package instance

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	OpStart   = iota // 0
	OpResume         // 1
	OpRestart        // 2
)

// RunOptions the options when running a FlowAction
type RunOptions struct {
	Op           int
	ReturnID     bool
	FlowURI      string
	InitialState *Instance
	ExecOptions  *ExecOptions
}

// ExecOptions are optional Patch & Interceptor to be used during instance execution
type ExecOptions struct {
	Patch       *support.Patch
	Interceptor *support.Interceptor
}

// IDGenerator generates IDs for flow instances
type IDGenerator interface {
	//NewFlowInstanceID generate a new instance ID
	NewFlowInstanceID() string
}

// ApplyExecOptions applies any execution options to the flow instance
func ApplyExecOptions(instance *Instance, execOptions *ExecOptions) {

	if execOptions != nil {

		if execOptions.Patch != nil {
			logger.Infof("Instance [%s] has patch", instance.ID())
			instance.Patch = execOptions.Patch
			instance.Patch.Init()
		}

		if execOptions.Interceptor != nil {
			logger.Infof("Instance [%s] has interceptor", instance.ID())
			instance.Interceptor = execOptions.Interceptor
			instance.Interceptor.Init()
		}
	}
}

// IDResponse is a response object consists of an ID
type IDResponse struct {
	ID string `json:"id"`
}
