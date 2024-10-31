package views

import (
	"horsaen/afreeca-downloader/tools"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	keyword = tools.MakeFgStyle("211")
	subtle  = tools.MakeFgStyle("241")
	dot     = tools.ColorFg(" • ", "236")
)

type (
	errMsg error
)

type model struct {
	Platform         int
	PlatformSelected bool
	Mode             int
	ModeSelected     bool
	UserSelect       bool
	TextInput        textinput.Model
	TextArea         textarea.Model
	TextLines        []string
	Running          bool
	err              error
}

func (m model) Init() tea.Cmd {
	return nil
}

func InitialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 156

	ta := textarea.New()
	ta.Focus()
	ta.CharLimit = 0

	return model{
		Platform:         0,
		PlatformSelected: false,
		Mode:             0,
		ModeSelected:     false,
		TextInput:        ti,
		TextArea:         ta,
		TextLines:        []string{},
		Running:          false,
		err:              nil,
	}
}

// this entire TUI could be cleaned up a bit, but whatever it looks good enough
func (m model) View() string {
	var s string

	if !m.Running {
		if !m.PlatformSelected {
			s = Platforms(m)
		} else if !m.ModeSelected {
			s = Modes(m)
		}
		// else if m.ModeSelected {
		// 	s = Afreeca(m)
		// }
	}

	return s
}
