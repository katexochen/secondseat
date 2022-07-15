package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/katexochen/secondseat/internal/xinput"
	"github.com/spf13/cobra"
)

const masterName = "secondseat"

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add input devices for a second user.",
		Args:  cobra.NoArgs,
		RunE:  runAdd,
	}
	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	xiHandler, err := xinput.NewHandler()
	if err != nil {
		return err
	}
	return add(cmd, xiHandler)
}

func add(cmd *cobra.Command, xiHandler xinput.Handler) error {
	pointerMaster, keyboardMaster, err := xiHandler.CreatePrimary(masterName)
	if err != nil {
		return err
	}
	fmt.Println(keyboardMaster)
	fmt.Printf("ID: %d, mID: %d\n", keyboardMaster.ID, keyboardMaster.PrimaryID)
	fmt.Println(pointerMaster)
	fmt.Printf("ID: %d, mID: %d\n", pointerMaster.ID, pointerMaster.PrimaryID)
	if err := addInput(cmd, xiHandler, pointerMaster.ID, xinput.Pointer); err != nil {
		return err
	}
	if err := addInput(cmd, xiHandler, keyboardMaster.ID, xinput.Keyboard); err != nil {
		return err
	}
	cmd.Println("Have fun together! (★•ヮ•)八(•ヮ•)")
	return nil
}

func addInput(cmd *cobra.Command, xiHandler xinput.Handler, masterID int, deviceType xinput.DeviceType) error {
	cmd.Printf("Make sure the %s for the second seat is disconnected.", deviceType)
	if err := confirmWithEnter(cmd.InOrStdin()); err != nil {
		return err
	}
	if err := xiHandler.UpdateState(); err != nil {
		return err
	}
	cmd.Printf("Now connect the second %s.", deviceType)
	if err := confirmWithEnter(cmd.InOrStdin()); err != nil {
		return err
	}
	newInputs, err := xiHandler.DetectNewSecondaries(deviceType)
	if err != nil {
		return err
	}
	cmd.Println("The following new devices were detected:")
	printInputs(cmd, newInputs)
	cmd.Printf("These devices will be attached to second seat.")
	if err := confirmWithEnter(cmd.InOrStdin()); err != nil {
		return err
	}
	if err := xiHandler.Reattach(newInputs, masterID); err != nil {
		return err
	}
	cmd.Println("Successfully reattached!")
	return nil
}

func printInputs(cmd *cobra.Command, inputs []xinput.Xinput) {
	for _, i := range inputs {
		cmd.Printf("  %s\n", i.Name)
	}
}

func confirmWithEnter(r io.Reader) error {
	reader := bufio.NewReader(r)
	_, err := reader.ReadString('\n')
	return err
}
