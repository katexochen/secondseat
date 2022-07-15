package xinput

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

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

func (h *Handler) GetPrimaries() []Xinput {
	stateSlice := stateToSlice(h.state)
	return filterInputs(stateSlice, Primary, "")
}

func (h *Handler) GetPrimariesByName(name string) (Xinput, Xinput, error) {
	primaries := h.GetPrimaries()
	var matches []Xinput
	for _, m := range primaries {
		if strings.Contains(m.Name, name) {
			matches = append(matches, m)
		}
	}
	return checkPrimary(matches)
}

func (h *Handler) Reattach(inputs []Xinput, primaryID int) error {
	for _, i := range inputs {
		if err := h.reattach(i.ID, primaryID); err != nil {
			return err
		}
	}
	return h.UpdateState()
}

func (h *Handler) CreatePrimary(name string) (Xinput, Xinput, error) {
	if err := h.createPrimary(name); err != nil {
		return Xinput{}, Xinput{}, err
	}
	newPrimaries, err := h.DetectNewPrimaries()
	if err != nil {
		return Xinput{}, Xinput{}, nil
	}
	return checkPrimary(newPrimaries)
}

func (h *Handler) RemovePrimary(id int) error {
	if err := h.removePrimary(id); err != nil {
		return err
	}
	return h.UpdateState()
}

func (h *Handler) DetectNewPrimaries() ([]Xinput, error) {
	newEntries, err := h.detectNew()
	if err != nil {
		return nil, err
	}
	return filterInputs(newEntries, Primary, ""), nil
}

func (h *Handler) DetectNewSecondaries(dType DeviceType) ([]Xinput, error) {
	newEntries, err := h.detectNew()
	if err != nil {
		return nil, err
	}
	return filterInputs(newEntries, Secondary, dType), nil
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

func checkPrimary(primary []Xinput) (Xinput, Xinput, error) {
	if len(primary) != 2 {
		return Xinput{}, Xinput{}, fmt.Errorf("expected 2 new primary inputs, got %d", len(primary))
	}

	keyboardPrimary := primary[0]
	pointerPrimary := primary[1]

	if keyboardPrimary.Type != Keyboard {
		keyboardPrimary, pointerPrimary = pointerPrimary, keyboardPrimary
	}
	if keyboardPrimary.Type != Keyboard ||
		pointerPrimary.Type != Pointer {
		return Xinput{}, Xinput{}, errors.New("new primary inputs have invalid types")
	}
	if keyboardPrimary.ID != pointerPrimary.PrimaryID ||
		pointerPrimary.ID != keyboardPrimary.PrimaryID {
		return Xinput{}, Xinput{}, errors.New("new primary inputs do not point to each other")
	}
	return pointerPrimary, keyboardPrimary, nil
}

func stateToSlice(state map[int]Xinput) []Xinput {
	v := make([]Xinput, 0, len(state))
	for _, value := range state {
		v = append(v, value)
	}
	return v
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
