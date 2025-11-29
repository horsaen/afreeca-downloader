package views

import (
	"horsaen/afreeca-downloader/plugins/bigo"
	"horsaen/afreeca-downloader/plugins/chzzk"
	"horsaen/afreeca-downloader/plugins/flex"
	"horsaen/afreeca-downloader/plugins/panda"
	"horsaen/afreeca-downloader/plugins/soop"
	"horsaen/afreeca-downloader/plugins/tiktok"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.PlatformSelected || m.ModeSelected && m.Mode == 0 {
		m.TextInput, _ = m.TextInput.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.Running && (strings.TrimSpace(m.TextInput.Value()) != "" || m.Platform == 8) && m.Mode != 2 {
				m.Running = true
				switch m.Platform {
				case 0:
					// handle soop modes
					switch m.Mode {
					case 0:
						go soop.Start(strings.TrimSpace(m.TextInput.Value()))
					case 1:
						go func() {
							soop.DownloadVod(strings.TrimSpace(m.TextInput.Value()))
							m.Running = false
							os.Exit(0)
						}()
					}
				case 1:
					go bigo.Start(strings.TrimSpace(m.TextInput.Value()))
				case 2:
					go chzzk.Start(strings.TrimSpace(m.TextInput.Value()))
				case 3:
					go flex.Start(strings.TrimSpace(m.TextInput.Value()))
				case 4:
					go panda.Start(strings.TrimSpace(m.TextInput.Value()))
				case 5:
					go tiktok.Start(strings.TrimSpace(m.TextInput.Value()))
					// case 6:
					// 	tools.ClearCli()
					// go concurrent.Start()
				}
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	if !m.PlatformSelected {
		return UpdatePlatforms(msg, m)
	}

	if m.PlatformSelected && !m.ModeSelected {
		return UpdateMode(msg, m)
	}

	return m, cmd
}

func UpdatePlatforms(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Platform++
			if m.Platform > 5 {
				m.Platform = 5
			}
		case "k", "up":
			m.Platform--
			if m.Platform < 0 {
				m.Platform = 0
			}
		case "enter":
			m.PlatformSelected = true
		}
	}

	return m, nil
}

func UpdateMode(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Mode++
			if m.Mode > 2 {
				m.Mode = 2
			}
		case "k", "up":
			m.Mode--
			if m.Mode < 0 {
				m.Mode = 0
			}
		case "enter":
			m.ModeSelected = true
		}
	}

	return m, nil
}
