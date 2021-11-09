package rpc

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/rs/xid"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	"github.com/tinkerbell/pbnj/pkg/logging"
	"github.com/tinkerbell/pbnj/pkg/task"
	"os/exec"
	"time"
)

// MachineService for doing power and device actions.
type MachineServiceGovc struct {
	Log logging.Logger
	// Timeout is how long a task should be run
	// before it is cancelled. This is for use in a
	// TaskRunner.Execute function that runs all BMC
	// interactions in the background.
	Timeout    time.Duration
	TaskRunner task.Task
	v1.UnimplementedMachineServer
}


func (m *MachineServiceGovc) Power(ctx context.Context, in *v1.PowerRequest) (*v1.PowerResponse, error) {
	l := m.Log.GetContextLogger(ctx)
	taskID := xid.New().String()
	l = l.WithValues("taskID", taskID)
	l.Info(
		"start Power request",
		"username", in.Authn.GetDirectAuthn().GetUsername(),
		"vendor", in.Vendor.GetName(),
		"powerAction", in.GetPowerAction().String(),
		"softTimeout", in.SoftTimeout,
		"OffDuration", in.OffDuration,
	)

	host := in.Authn.GetDirectAuthn().GetHost().Host
	//userName := in.Authn.GetDirectAuthn().GetUsername()
	//password := in.Authn.GetDirectAuthn().GetUsername()
	app := "govc"

	arg0 := "vm.power"
	arg1 := "-on"
	arg2 := host
	print(app+" "+arg0+" "+arg1+" "+arg2)
	cmd := exec.Command(app, arg0, arg1, arg2)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return &v1.PowerResponse{TaskId: taskID}, nil
	}

	// Print the output
	fmt.Println(string(stdout))

	return &v1.PowerResponse{TaskId: taskID}, nil
}