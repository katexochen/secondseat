package cmd

import (
	"github.com/katexochen/secondseat/internal/xinput"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove second input user",
		Args:  cobra.NoArgs,
		RunE:  runRemove,
	}
	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	handler, err := xinput.NewHandler()
	if err != nil {
		return err
	}

	err = remove(cmd, handler)
	if err != nil {
		return err
	}

	cmd.Println("Second seat successfully removed.")
	return nil
}

func remove(cmd *cobra.Command, handler xinput.Handler) error {
	pointerMaster, _, err := handler.GetPrimariesByName(masterName)
	if err != nil {
		return err
	}
	return handler.RemovePrimary(pointerMaster.ID)
}
