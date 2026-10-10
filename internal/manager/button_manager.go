package manager

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/IridiumNan/link-manager/internal/models"
)

type ButtonMgr struct {
	buttons []models.NameButtonOperator

	// store the current pos
	curr int

	// verbose flag
	verbose bool
}

// TurnOn current button
func (m *ButtonMgr) TurnOn() {
	m.buttons[m.curr].TurnOn()
}

// TurnOff current button
func (m *ButtonMgr) TurnOff() {
	m.buttons[m.curr].TurnOff()
}

// Switch current button
func (m *ButtonMgr) Switch() {
	m.buttons[m.curr].Switch()
}

// Status return current button status
func (m *ButtonMgr) Status() models.Status {
	return m.buttons[m.curr].Status()
}

// Prev move the button pointer to previous one
// if idx < 0, it will set it to 0
func (m *ButtonMgr) Prev() {
	m.curr--

	if m.curr < 0 {
		m.curr = 0
	}
}

// Next move the button pointer to next one
// if idx out of range, it will set it to max valid idx
func (m *ButtonMgr) Next() {
	m.curr++

	maxIdx := len(m.buttons) - 1

	if m.curr > maxIdx {
		m.curr = maxIdx
	}
}

// CheckBox return a string with format
// [ ] buttons[curr] for button with [StatusOff]
// [*] buttons[curr] for button with [StatusOn]
// [?] buttons[curr] for button with [StatusNone]
// end with a new line '\n'
//
// if curr out of range, it will return empty string
func (m *ButtonMgr) CheckBox() string {
	return m.checkBox(m.curr)
}

func (m *ButtonMgr) checkBox(idx int) string {
	if idx >= len(m.buttons) || idx < 0 {
		return ""
	}

	b := m.buttons[idx]
	// default mark is [StatusNone]
	mark := "?"
	switch b.Status() {
	case models.StatusOff:
		mark = " "
	case models.StatusOn:
		mark = "*"
	}

	return fmt.Sprintf("[%v] %s\n", mark, b.String())
}

// String is the implemnt for [models.Name]
// return all [ButtonMgr.checkBox]
func (m *ButtonMgr) String() string {
	var s strings.Builder

	s.Grow(1 << 8)

	for i := range m.buttons {
		s.WriteString(m.checkBox(i))
	}

	return s.String()
}

// Sync change current button status
func (m *ButtonMgr) Sync() (changeMsg string, errs []error) {
	errs = make([]error, 0, 5)

	var msgCollector strings.Builder
	msgCollector.Grow(1 << 6)

	for _, b := range m.buttons {
		msg, err := b.Operate()
		if err != nil {
			errs = append(errs, err)

			_, err := fmt.Fprintf(&msgCollector, "%s\n", string(msg))
			if err != nil {
				slog.Error("error when write message into string builder", "err", err)
			}
			continue
		}

		if !m.verbose && (msg == models.MsgNoChange) {
			continue
		}

		_, err = fmt.Fprintf(&msgCollector, "%s\n", string(msg))
		if err != nil {
			slog.Error("error when write message into string builder", "err", err)
		}
	}

	changeMsg = msgCollector.String()

	if len(errs) == 0 {
		errs = nil
		return
	}

	return
}

// NewButtonMgr create a new manager with given buttons
// set the verbose msg collection with [SetVerbose](true), default false
func NewButtonMgr(buttons []models.NameButtonOperator) *ButtonMgr {
	return &ButtonMgr{
		buttons: buttons,
		curr:    0,
		verbose: false,
	}
}

// SetVerbose the verbose msg collection
func (m *ButtonMgr) SetVerbose(verbose bool) {
	m.verbose = verbose
}
