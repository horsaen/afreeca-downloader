package views

import (
	"fmt"
	"horsaen/afreeca-downloader/tools"
)

func Platforms(m model) string {
	c := m.Platform

	tpl := "Select Platform:\n\n"
	tpl += "%s\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	Platforms := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		tools.Checkbox("Afreeca", c == 0),
		tools.Checkbox("Bigo", c == 1),
		tools.Checkbox("Chzzk", c == 2),
		tools.Checkbox("Flex", c == 3),
		tools.Checkbox("Kick", c == 4),
		tools.Checkbox("Panda", c == 5),
		tools.Checkbox("TikTok", c == 6),
		tools.Checkbox("Soop", c == 7),
		tools.Checkbox("Concurrent", c == 8),
	)

	return fmt.Sprintf(tpl, Platforms)
}

func Modes(m model) string {
	var msg string

	switch m.Platform {
	case 0:
		msg = AfreecaModes(m)
	case 1:
		// bigo
		msg = fmt.Sprintf(
			"Username\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case 2:
		// chzzk
		msg = fmt.Sprintf(
			"Username\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case 3:
		// flex
		msg = fmt.Sprintf(
			"User ID:\n%s\n\n%s\n\n%s",
			"https://www.flextv.co.kr/channels/"+keyword("xxxxxx")+"/live",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case 4:
		// kick
		msg = fmt.Sprintf(
			"Username:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case 5:
		// panda
		msg = fmt.Sprintf(
			"Username:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case 6:
		// tiktok
		msg = fmt.Sprintf(
			"Username:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case 7:
		// soop
		msg = fmt.Sprintf(
			"Username:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"

	}

	return msg
}

func AfreecaModes(m model) string {
	c := m.Mode

	tpl := "Select Mode:\n\n"
	tpl += "%s\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	Modes := fmt.Sprintf(
		"%s\n%s\n%s\n",
		tools.Checkbox("Stream", c == 0),
		tools.Checkbox("VOD", c == 1),
		tools.Checkbox("Playlist", c == 2),
	)

	return fmt.Sprintf(tpl, Modes)
}

func Afreeca(m model) string {
	var msg string

	switch m.Mode {
	case 0:
		msg = fmt.Sprintf(
			"Username:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case 1:
		msg = fmt.Sprintf(
			"VOD ID:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case 2:
		msg = fmt.Sprintf(
			"One m3u8 per line.\n\n%s\n\n%s",
			m.TextArea.View(),
			"(ctrl+c/esc to start download)",
		) + "\n\n"
	}

	return msg
}
