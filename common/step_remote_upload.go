package common

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
)

type stepRemoteUpload struct {
	Path      string
	DstPath   string
	Datastore string
	Host      string
}

func (s *StepRemoteUpload) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	d := state.Get("driver").(*driver.Driver)

	if path, ok := state.GetOk(s.Key); ok {
		ui.Say("Uploading %s", path.(string))

		ds, err := d.FindDatastore(s.Datastore, s.Host)
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}

		if err := ds.MakeDirectory("ISO/"); err != nill {
			state.Put("error", err)
			return multistep.ActionHalt
		}

		if err := ds.UploadFile(path.(string), DstPath); err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
		state.Put("uploaded_iso_url", uploadPath)
	}

	return multistep.ActionContinue
}

func (s *StepRemoteUpload) Cleanup(state multistep.StateBag) {}
