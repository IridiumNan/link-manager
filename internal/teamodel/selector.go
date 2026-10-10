package teamodel

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/fatih/color"
)

var selectedColor = color.BgRGB(110, 110, 110)

// Selector is implement for [tea.Model]
type Selector struct {
	// header will be added in the [Selector.View] str head
	header string

	// footer will be append in tail of [Selector.View]
	footer string

	choices []string

	cursor int
}

// Init empty implement for [tea.Model]
func (s *Selector) Init() tea.Cmd {
	return nil
}

// Update handle the key press msg
func (s *Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit

		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}

		case "down", "j":
			if s.cursor < len(s.choices)-1 {
				s.cursor++
			}

		case "enter":
			return s, tea.Quit
		}
	}

	return s, nil
}

// indexView return a string view specific choice
// it used for [Selector.View] function
func (s *Selector) idxView(idx int) string {
	format := "%s %s\n"
	cursor := " "

	if idx == s.cursor {
		cursor = ">"

		return selectedColor.Sprintf(format, cursor, s.choices[idx])
	}

	return fmt.Sprintf(format, cursor, s.choices[idx])
}

// View is implement of [tea.Model]
// Write it will formated as
//
// [Selector.header]
// [Selector.idxView] for each choice
// [Selector.footer]
func (s *Selector) View() tea.View {
	var str strings.Builder
	str.Grow(1 << 8)

	str.WriteString(s.header)

	for i := range s.choices {
		str.WriteString(s.idxView(i))
	}

	str.WriteString(s.footer)

	return tea.NewView(str.String())
}

// Choice return current choice
// you should only call it when ended
func (s *Selector) Choice() string {
	return s.choices[s.cursor]
}

// NewSelector chreate new selector
func NewSelector(choices []string) *Selector {
	if len(choices) == 0 {
		return nil
	}

	return &Selector{
		choices: choices,
		cursor:  0,
	}
}

// SetHeader set the header which will be used on [Selector.View]
func (s *Selector) SetHeader(header string) {
	s.header = header
}

// SetFooter set the footer which will be used on [Selector.View]
func (s *Selector) SetFooter(footer string) {
	s.footer = footer
}
