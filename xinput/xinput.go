package xinput

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/hashicorp/go-version"
)

const requiredXVersion = "1.20"

type Xinput struct {
	Name string
	ID   int
}

func List() ([]Xinput, error) {
	cmd := exec.Command("xinput", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	reg, err := regexp.Compile(`(?:↳ )(\S+(?: \S+)*)|(?:id=)(\d+)(?:\s\[slave)`)
	if err != nil {
		return nil, err
	}

	matches := reg.FindAllStringSubmatch(string(out), -1)

	var result []Xinput
	var nextInput Xinput
	for _, match := range matches {
		switch {
		case len(match[1]) != 0:
			nextInput.Name = match[1]

		case len(match[2]) != 0:
			id, err := strconv.Atoi(match[2])
			if err != nil {
				return nil, err
			}
			nextInput.ID = id
			if len(nextInput.Name) == 0 {
				return nil, fmt.Errorf("could not find a name to input id %d", id)
			}
			result = append(result, nextInput)
			nextInput = Xinput{}
		default:
			return nil, errors.New("got regex match with two empty capturing groups")
		}
	}
	return result, nil
}

func Reattach(inputID, masterID int) error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "reattach")
	cmdArgs = append(cmdArgs, strconv.Itoa(inputID))
	cmdArgs = append(cmdArgs, strconv.Itoa(masterID))
	cmd := exec.Command("xinput", cmdArgs...)
	_, err := cmd.Output()
	return err
}

func CreateMaster(name string) error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "create-master")
	cmdArgs = append(cmdArgs, name)
	cmd := exec.Command("xinput", cmdArgs...)
	_, err := cmd.Output()
	return err
}

func CheckXServerVersion() (bool, error) {
	cmd := exec.Command("dpkg", "-l")
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	reg, err := regexp.Compile(`(?:xserver-xorg-core[^:]*:)([\d\.]*)`)
	if err != nil {
		return false, err
	}

	foundStr := reg.FindStringSubmatch(string(out))[1]
	found, err := version.NewVersion(foundStr)
	if err != nil {
		return false, err
	}
	required, err := version.NewVersion(requiredXVersion)
	if err != nil {
		return false, err
	}
	return found.GreaterThanOrEqual(required), nil
}
