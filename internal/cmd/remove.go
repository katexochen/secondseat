package cmd

import (
	"github.com/katexochen/secondseat/internal/xinput"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove second input user.",
		Long:  "Remove second input user",
		Args:  cobra.NoArgs,
		RunE:  runRemove,
	}
	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	xiHandler, err := xinput.NewHandler()
	if err != nil {
		return err
	}
	return remove(cmd, xiHandler)
}

func remove(cmd *cobra.Command, xiHandler xinput.Handler) error {
	pointerMaster, _, err := xiHandler.GetMasterByName(masterName)
	if err != nil {
		return err
	}
	return xiHandler.RemoveMaster(pointerMaster.ID)
}
