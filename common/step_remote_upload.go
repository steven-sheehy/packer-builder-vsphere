package common

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"path/filepath"
)

type StepRemoteUpload struct {
	Path      string
	DstPath   string
	Datastore string
	Host      string
}

func (s *StepRemoteUpload) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	d := state.Get("driver").(*driver.Driver)

	if path, ok := state.GetOk(s.Path); ok {
		filename := filepath.Base(path.(string))
		remotepath := fmt.Sprintf("ISO/%s", filename)
		remotedirectory := fmt.Sprintf("[%s] ISO/", s.Datastore)
		final_remotepath := fmt.Sprintf("[%s] %s", s.Datastore, remotepath)

		ui.Say(fmt.Sprintf("Uploading %s to %s", filename, remotepath))

		ds, err := d.FindDatastore(s.Datastore, s.Host)
		if err != nil {
			ui.Say("Datastore doesn't exist")
			state.Put("error", err)
			return multistep.ActionHalt
		}

		if exists := ds.FileExists(remotepath); exists == true {
			ui.Say("File already upload")
			state.Put(s.DstPath, final_remotepath)
			state.Put("error", err)
			return multistep.ActionContinue
		}

		if err := ds.MakeDirectory(remotedirectory); err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}

		if err := ds.UploadFile(path.(string), remotepath); err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
		state.Put(s.DstPath, final_remotepath)
	}

	return multistep.ActionContinue
}

func (s *StepRemoteUpload) Cleanup(state multistep.StateBag) {}
