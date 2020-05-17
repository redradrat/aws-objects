package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/redradrat/cloud-objects/cloudobject"
)

const (
	CreateCloudObjectAction CloudObjectAction = "create"
	ReadCloudObjectAction   CloudObjectAction = "read"
	UpdateCloudObjectAction CloudObjectAction = "update"
	DeleteCloudObjectAction CloudObjectAction = "delete"
)

type CloudObjectAction string

func OnlyCloudObjectAction() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("accepts %d arg(s), received %d", 1, len(args))
		}
		arg := args[0]
		if !IsCloudObjectAction(arg) {
			return fmt.Errorf("unknown action argument received %s", arg)
		}
		return nil
	}
}

func IsCloudObjectAction(s string) bool {
	switch CloudObjectAction(s) {
	case CreateCloudObjectAction:
		return true
	case ReadCloudObjectAction:
		return true
	case UpdateCloudObjectAction:
		return true
	case DeleteCloudObjectAction:
		return true
	default:
		return false
	}
}

func HandleCloudObject(obj cloudobject.CloudObject, spec cloudobject.CloudObjectSpec,
	action CloudObjectAction, purge bool) (cloudobject.Secrets, error) {
	switch action {
	case CreateCloudObjectAction:
		return obj.Create(spec)
	case ReadCloudObjectAction:
		return nil, obj.Read()
	case UpdateCloudObjectAction:
		return obj.Update(spec)
	case DeleteCloudObjectAction:
		return nil, obj.Delete(purge)
	default:
		return nil, CloudObjectActionUnknown{Message: fmt.Sprintf("action '%s' unknown", string(action))}
	}
}

type CloudObjectActionUnknown struct {
	Message string
}

func (e CloudObjectActionUnknown) Error() string {
	return e.Message
}
