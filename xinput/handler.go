package xinput

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/hashicorp/go-version"
)

type Handler struct {
	state        map[int]Xinput
	newCommander func() Commander
}

func NewHandler() (Handler, error) {
	if err := checkXServerVersion(); err != nil {
		return Handler{}, err
	}
	h := Handler{newCommander: NewExecCommander}
	if err := h.UpdateState(); err != nil {
		return Handler{}, err
	}
	return h, nil
}

// for Debugging
func (h *Handler) GetState() map[int]Xinput {
	return h.state
}

func (h *Handler) Reattach(inputs []Xinput, masterID int) error {
	for _, i := range inputs {
		if err := h.reattach(i.ID, masterID); err != nil {
			return err
		}
	}
	return h.UpdateState()
}

func (h *Handler) CreateMaster(name string) (Xinput, Xinput, error) {
	if err := h.createMaster(name); err != nil {
		return Xinput{}, Xinput{}, err
	}
	newMasters, err := h.DetectNewMasters()
	if err != nil {
		return Xinput{}, Xinput{}, nil
	}
	if len(newMasters) != 2 {
		return Xinput{}, Xinput{}, fmt.Errorf("expected 2 new master, got %d", len(newMasters))
	}

	keyboardMaster := newMasters[0]
	pointerMaster := newMasters[1]

	if keyboardMaster.Type != Keyboard {
		keyboardMaster, pointerMaster = pointerMaster, keyboardMaster
	}
	if keyboardMaster.Type != Keyboard ||
		pointerMaster.Type != Pointer {
		return Xinput{}, Xinput{}, errors.New("new masters have invalid types")
	}
	if keyboardMaster.ID != pointerMaster.MasterID ||
		pointerMaster.ID != keyboardMaster.MasterID {
		return Xinput{}, Xinput{}, errors.New("new masters do not point to each other")
	}
	return keyboardMaster, pointerMaster, nil
}

func (h *Handler) RemoveMaster(id int) error {
	if err := h.removeMaster(id); err != nil {
		return err
	}
	return h.UpdateState()
}

func (h *Handler) DetectNewMasters() ([]Xinput, error) {
	newEntries, err := h.detectNew()
	if err != nil {
		return nil, err
	}
	return filterInputs(newEntries, Master, ""), nil
}

func (h *Handler) DetectNewSlaves(dType DeviceType) ([]Xinput, error) {
	newEntries, err := h.detectNew()
	if err != nil {
		return nil, err
	}
	return filterInputs(newEntries, Slave, dType), nil
}

func (h *Handler) UpdateState() error {
	state, err := h.list()
	if err != nil {
		return err
	}
	h.state = state
	return nil
}

func (h *Handler) detectNew() ([]Xinput, error) {
	oldState := h.state
	if err := h.UpdateState(); err != nil {
		return nil, err
	}
	var newInputs []Xinput
	for k, v := range h.state {
		if _, found := oldState[k]; !found {
			newInputs = append(newInputs, v)
		}
	}
	return newInputs, nil
}

func filterInputs(inputs []Xinput, dRole DeviceRole, dType DeviceType) []Xinput {
	var matches []Xinput
	for _, i := range inputs {
		if i.Role != dRole && dRole != "" {
			continue
		}
		if i.Type != dType && dType != "" {
			continue
		}
		matches = append(matches, i)
	}
	return matches
}

func checkXServerVersion() error {
	cmd := exec.Command("dpkg", "-l")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	reg, err := regexp.Compile(`(?:xserver-xorg-core[^:]*:)([\d\.]*)`)
	if err != nil {
		return err
	}
	foundStr := reg.FindStringSubmatch(string(out))[1]
	found, err := version.NewVersion(foundStr)
	if err != nil {
		return err
	}
	required, err := version.NewVersion(requiredXVersion)
	if err != nil {
		return err
	}
	if found.LessThan(required) {
		return fmt.Errorf("installed X version %s is lower than required version %s", found.String(), required.String())
	}
	return nil
}
