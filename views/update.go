package views

import (
	"horsaen/afreeca-downloader/plugins/afreeca"
	"horsaen/afreeca-downloader/plugins/bigo"
	"horsaen/afreeca-downloader/plugins/chzzk"
	"horsaen/afreeca-downloader/plugins/concurrent"
	"horsaen/afreeca-downloader/plugins/flex"
	"horsaen/afreeca-downloader/plugins/kick"
	"horsaen/afreeca-downloader/plugins/panda"
	"horsaen/afreeca-downloader/plugins/soop"
	"horsaen/afreeca-downloader/plugins/tiktok"
	"horsaen/afreeca-downloader/tools"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.PlatformSelected || m.ModeSelected && m.Mode == 0 {
		m.TextInput, _ = m.TextInput.Update(msg)
	}

	if m.ModeSelected && m.Mode == 2 {
		m.TextArea, cmd = m.TextArea.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.Mode == 2 && !m.Running {
				m.Running = true
				done := make(chan bool)

				go func() {
					afreeca.Playlist(strings.Split(m.TextArea.Value(), "\n"))
					done <- true
				}()

				go func() {
					<-done
					m.Running = false
					os.Exit(0)
				}()

			} else {
				return m, tea.Quit
			}
		case tea.KeyEnter:
			if !m.Running && (m.TextInput.Value() != "" || m.Platform == 6) && m.Mode != 2 {
				m.Running = true
				switch m.Platform {
				case 0:
					switch m.Mode {
					case 0:
						go afreeca.Start(m.TextInput.Value())
					case 1:
						done := make(chan bool)

						go func() {
							afreeca.Vod(m.TextInput.Value())
							done <- true
						}()

						go func() {
							<-done
							m.Running = false
							os.Exit(0)
						}()
					}
				case 1:
					go bigo.Start(m.TextInput.Value())
				case 2:
					go chzzk.Start(m.TextInput.Value())
				case 3:
					go flex.Start(m.TextInput.Value())
				case 4:
					go kick.Start(m.TextInput.Value())
				case 5:
					go panda.Start(m.TextInput.Value())
				case 6:
					go tiktok.Start(m.TextInput.Value())
				case 7:
					go soop.Start(m.TextInput.Value())
				case 8:
					tools.ClearCli()
					go concurrent.Start()
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
			if m.Platform > 8 {
				m.Platform = 8
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
