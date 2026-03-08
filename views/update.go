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

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.PlatformSelected && m.ModeSelected && !m.LoginPrompt {
		m.TextInput, _ = m.TextInput.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.LoginPrompt {
			switch msg.String() {
			case "y", "Y":
				m.LoginPrompt = false
				m.NeedsLogin = true
				m.TextInput.SetValue("")
				return m, nil
			case "n", "N", "esc":
				m.LoginPrompt = false
				m.NeedsLogin = false
				m.TextInput.SetValue("")
				return m, nil
			}
			return m, nil
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.ModeSelected && m.Platform == 0 && m.Mode == 0 && m.NeedsLogin {
				switch m.LoginStep {
				case 0:
					m.LoginUser = strings.TrimSpace(m.TextInput.Value())
					m.TextInput.SetValue("")
					m.LoginStep = 1
					m.TextInput.EchoMode = textinput.EchoPassword
					return m, nil
				case 1:
					m.LoginPass = m.TextInput.Value()
					m.TextInput.SetValue("")
					m.TextInput.EchoMode = textinput.EchoNormal
					m.NeedsLogin = false
					soop.InitCredentials(m.LoginUser, m.LoginPass)
					return m, nil
				}
			}

			if !m.Running && (strings.TrimSpace(m.TextInput.Value()) != "" || m.Platform == 8) {
				m.Running = true
				switch m.Platform {
				case 0:
					switch m.Mode {
					case 0:
						go soop.Start(strings.TrimSpace(m.TextInput.Value()))
					case 1:
						go func() {
							soop.DownloadPlaylist(strings.TrimSpace(m.TextInput.Value()))
							m.Running = false
							os.Exit(0)
						}()
					case 2:
						go func() {
							soop.UserVods(strings.TrimSpace(m.TextInput.Value()))
							m.Running = false
							os.Exit(0)
						}()
					case 3:
						go func() {
							soop.DownloadUserVod(strings.TrimSpace(m.TextInput.Value()))
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
			if m.Mode > 3 {
				m.Mode = 3
			}
		case "k", "up":
			m.Mode--
			if m.Mode < 0 {
				m.Mode = 0
			}
		case "enter":
			m.ModeSelected = true
			if m.Platform == 0 && m.Mode == 0 {
				m.LoginPrompt = true
				m.TextInput.SetValue("")
			}
		}
	}

	return m, nil
}
