package teamodel

import (
	"fmt"
	"log/slog"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/IridiumNan/link-manager/internal/manager"
	"github.com/IridiumNan/link-manager/internal/models"
	"github.com/fatih/color"
)

// LinkManager is the implement of [tea.Model]
// who wrap a [manager.ButtonMgr] then update status
type LinkManager struct {
	header string

	footer string

	mgr manager.ButtonMgr
}

func (lm *LinkManager) Init() tea.Cmd {
	return nil
}

// done call the manager Sync function then print output
func (lm *LinkManager) done() {
	cm, errs := lm.mgr.Sync()

	fmt.Println(cm)

	if len(errs) > 0 {
		color.Red("some errors happened on process, see log by: %s log", models.AppName)
	}

	for _, e := range errs {
		slog.Error(e.Error())
	}
}

// Update listening the key event then update ButtonMgr status
func (lm *LinkManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return lm, tea.Quit

		case "up", "k":
			lm.mgr.Prev()

		case "down", "j":
			lm.mgr.Next()

		case "space":
			lm.mgr.Switch()

		case "enter":
			lm.done()
			return lm, tea.Quit
		}
	}

	return lm, nil
}

// View display latest status of links by [manager.ButtonMgr.String]
func (lm *LinkManager) View() tea.View {
	var str strings.Builder

	str.Grow(1 << 8)

	str.WriteString(lm.header)
	str.WriteString(lm.mgr.String())
	str.WriteString(lm.footer)

	return tea.NewView(str.String())
}

// NewLinkManager create a new link manager
// Init with a button manager with multi buttons
//
// You can set the header by [LinkManager.SetHeader]
// and set footer by [LinkManager.SetFooter]
func NewLinkManager(mgr manager.ButtonMgr) *LinkManager {
	return &LinkManager{
		mgr: mgr,
	}
}

// SetHeader set the header which will be used on [LinkManager.View]
func (lm *LinkManager) SetHeader(header string) {
	lm.header = header
}

// SetFooter set the footer which will be used on [LinkManager.View]
func (lm *LinkManager) SetFooter(footer string) {
	lm.footer = footer
}
