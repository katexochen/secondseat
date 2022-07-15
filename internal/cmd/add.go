package cmd

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/katexochen/secondseat/internal/xinput"
	"github.com/spf13/cobra"
)

const (
	masterName             = "secondseat"
	sleepAfterInputConnect = 500 * time.Millisecond
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add input devices for a second user",
		Args:  cobra.NoArgs,
		RunE:  runAdd,
	}
	cmd.Flags().BoolP("debug", "d", false, "Enable debug mode")
	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	handler, err := xinput.NewHandler()
	if err != nil {
		return err
	}
	return add(cmd, handler)
}

func add(cmd *cobra.Command, handler xinput.Handler) error {
	debugMode, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}
	if debugMode {
		printDebug(cmd, "Debug mode enabled\n")
	}

	secondPointer, err := detectInput(cmd, handler, xinput.Pointer)
	if err != nil {
		return err
	}
	secondKeyboard, err := detectInput(cmd, handler, xinput.Keyboard)
	if err != nil {
		return err
	}

	newPrimaryPointer, newPrimaryKeyboard, err := handler.CreatePrimary(masterName)
	if err != nil {
		return err
	}
	cmd.Println("\nSuccessfully created a new primary device for second seat.")
	if debugMode {
		printDebug(cmd, "New PrimaryPointer created: %v\n", newPrimaryPointer)
		printDebug(cmd, "New PrimaryKeyboard created: %v\n", newPrimaryKeyboard)
	}

	if err := addInput(cmd, handler, newPrimaryPointer.ID, secondPointer, debugMode); err != nil {
		return err
	}
	if err := addInput(cmd, handler, newPrimaryKeyboard.ID, secondKeyboard, debugMode); err != nil {
		return err
	}
	cmd.Println("Successfully reattached devices to second seat.")

	cmd.Println("\nHave fun together! ( •ヮ•)八(•ヮ• )")
	return nil
}

func detectInput(cmd *cobra.Command, handler xinput.Handler, deviceType xinput.DeviceType) ([]xinput.Xinput, error) {
	cmd.Printf("\nMake sure the %s for the second seat is disconnected. [↵]", deviceType)
	if err := confirmWithEnter(cmd.InOrStdin()); err != nil {
		return nil, err
	}
	if err := handler.UpdateState(); err != nil {
		return nil, err
	}

	cmd.Printf("Now connect the second %s. [↵]", deviceType)
	if err := confirmWithEnter(cmd.InOrStdin()); err != nil {
		return nil, err
	}
	time.Sleep(sleepAfterInputConnect)
	newInputs, err := handler.DetectNewSecondaries(deviceType)
	if err != nil {
		return nil, err
	}
	if len(newInputs) == 0 {
		return nil, fmt.Errorf("no new %s detected", deviceType)
	}

	cmd.Printf("The following new %ss were detected:\n", deviceType)
	printInputs(cmd, newInputs)
	return newInputs, nil
}

func addInput(cmd *cobra.Command, handler xinput.Handler, primaryID int, inputs []xinput.Xinput, debug bool) error {
	if err := handler.Reattach(inputs, primaryID); err != nil {
		return err
	}
	if debug {
		printDebug(cmd, "Successfully reattached inputs %v to primary device %d\n", inputs, primaryID)
	}
	return nil
}

func confirmWithEnter(r io.Reader) error {
	reader := bufio.NewReader(r)
	_, err := reader.ReadString('\n')
	return err
}

func printInputs(cmd *cobra.Command, inputs []xinput.Xinput) {
	for _, i := range inputs {
		cmd.Printf("  %s\n", i.Name)
	}
}

func printDebug(cmd *cobra.Command, format string, args ...interface{}) {
	format = fmt.Sprintf("[DEBUG] %s", format)
	cmd.Printf(format, args...)
}
