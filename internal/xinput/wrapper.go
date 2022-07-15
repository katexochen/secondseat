package xinput

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// listRegx parses the output of xinput list.
// Group 1: Name
// Group 2: ID
// Group 3: DeviceRole
// Group 4: DeviceType
// Group 5: PrimaryID
var listRegx = regexp.MustCompile(`(?:[⎡⎜⎣↳ ]+)(\S+(?: \S+)*)(?:\s*id=)(\d+)(?:\s\[)(\w+)(?: +)(\w+)(?: +\()(\d+)`)

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

	matches := listRegx.FindAllStringSubmatch(string(out), -1)
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
		primaryID, err := strconv.Atoi(match[5])
		if err != nil {
			return nil, err
		}
		result[id] = Xinput{
			Name:      match[1],
			ID:        id,
			Type:      device,
			Role:      role,
			PrimaryID: primaryID,
		}
	}
	return result, nil
}

// reattach wraps 'xinput reattach' command.
func (h *Handler) reattach(inputID, primaryID int) error {
	fmt.Printf("reattaching %d to %d\n", inputID, primaryID)
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "reattach")
	cmdArgs = append(cmdArgs, strconv.Itoa(inputID))
	cmdArgs = append(cmdArgs, strconv.Itoa(primaryID))
	cmder := h.newCommander()
	cmder.Command("xinput", cmdArgs...)
	_, err := cmder.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Errorf("reattach failed with output '%s'", exitErr.Stderr)
	}
	return err
}

// createPrimary wraps 'xinput create-master' command.
func (h *Handler) createPrimary(name string) error {
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

// removePrimary wraps 'xinput remove-master' command.
//
// Removing a primary pointer will also remove the primary keyboard
// it is pointing to and viceversa.
func (h *Handler) removePrimary(id int) error {
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
