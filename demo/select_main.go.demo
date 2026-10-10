package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	choices []string

	cursor int

	selected map[int]bool

	saved bool
}

func NewModel() model {
	return model{
		choices: []string{
			"Fedora",
			"Ubuntu",
			"Arch",
			"Mint",
			"Nixos",
		},

		selected: make(map[int]bool),

		saved: false,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "space":
			// change the status directly
			m.selected[m.cursor] = !m.selected[m.cursor]

		case "enter":
			m.saved = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) View() tea.View {
	var s strings.Builder
	s.Grow(1 << 8)
	s.WriteString("Which distro should I choose ?\n\n")

	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "

		if m.selected[i] {
			checked = "x"
		}
		_, err := fmt.Fprintf(&s, "%s [%s] %s\n", cursor, checked, choice)
		if err != nil {
			slog.Error("error while write string into the view string content builder", "err", err)
		}
	}

	s.WriteString("\nPress q to quit\n")

	return tea.NewView(s.String())
}

func (m model) Fprint(out io.Writer) {
	var s strings.Builder

	s.Grow(2 << 6)

	s.WriteString("You have chosen the ")

	for i, choice := range m.choices {
		if m.selected[i] {
			fmt.Fprintf(&s, "%s ", choice)
		}
	}

	s.WriteString("distro")

	_, err := fmt.Fprint(out, s.String())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	m := NewModel()
	p := tea.NewProgram(&m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %s", err.Error())
		os.Exit(1)
	}

	if m.saved {
		m.Fprint(os.Stdout)
	}
}
