package xinput

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// list wraps 'xinput list' command
func (h *Handler) list() (map[int]Xinput, error) {
	cmder := h.newCommander()
	cmder.Command("xinput", "list")
	out, err := cmder.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return nil, fmt.Errorf("failed with output '%s'", exitErr.Stderr)
	}
	if err != nil {
		return nil, err
	}
	// Group 1: Name
	// Group 2: ID
	// Group 3: DeviceRole
	// Group 4: DeviceType
	// Group 5: MasterID
	reg, err := regexp.Compile(`(?:[⎡⎜⎣↳ ]+)(\S+(?: \S+)*)(?:\s*id=)(\d+)(?:\s\[)(\w+)(?: +)(\w+)(?: +\()(\d+)`)
	if err != nil {
		return nil, err
	}
	matches := reg.FindAllStringSubmatch(string(out), -1)

	result := make(map[int]Xinput)
	for _, match := range matches {

		for i, group := range match {
			if len(group) == 0 {
				return nil, fmt.Errorf("regex capturing group %d is empty", i)
			}
		}
		id, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}
		role, err := parseDeviceRole(match[3])
		if err != nil {
			return nil, err
		}
		device, err := parseDeviceType(match[4])
		if err != nil {
			return nil, err
		}
		masterID, err := strconv.Atoi(match[5])
		if err != nil {
			return nil, err
		}
		result[id] = Xinput{
			Name:     match[1],
			ID:       id,
			Type:     device,
			Role:     role,
			MasterID: masterID,
		}
	}
	return result, nil
}

// reattach wraps 'xinput reattach' command.
func (h *Handler) reattach(inputID, masterID int) error {
	fmt.Printf("reattaching %d to %d\n", inputID, masterID)
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "reattach")
	cmdArgs = append(cmdArgs, strconv.Itoa(inputID))
	cmdArgs = append(cmdArgs, strconv.Itoa(masterID))
	cmder := h.newCommander()
	cmder.Command("xinput", cmdArgs...)
	_, err := cmder.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Errorf("reattach failed with output '%s'", exitErr.Stderr)
	}
	return err
}

// createMaster wraps 'xinput create-master' command.
func (h *Handler) createMaster(name string) error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "create-master")
	cmdArgs = append(cmdArgs, name)
	cmder := h.newCommander()
	cmder.Command("xinput", cmdArgs...)
	_, err := cmder.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Errorf("create-master failed with output '%s'", exitErr.Stderr)
	}
	return err
}

// removeMaster wraps 'xinput remove-master' command.
//
// Removing a master pointer will also remove the master keyboard
// it is pointing to and viceversa.
func (h *Handler) removeMaster(id int) error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "remove-master")
	cmdArgs = append(cmdArgs, strconv.Itoa(id))
	cmder := h.newCommander()
	cmder.Command("xinput", cmdArgs...)
	_, err := cmder.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Errorf("remove-master failed with output '%s'", exitErr.Stderr)
	}
	return err
}
