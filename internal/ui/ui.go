package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hnpf/beetrap/internal/capture"
)

var (
	_ssh  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F5F")).Padding(0, 1)
	_ftp  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5F9FFF")).Padding(0, 1)
	_http = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5FFF87")).Padding(0, 1)
	_dim  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	_bold = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#EEEEEE"))
	_hdr  = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFB347")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF8C00")).
		Padding(0, 2)
)

type _connMsg capture.VxConnection

type VxModel struct {
	conns  []capture.VxConnection
	ch     chan capture.VxConnection
	counts map[string]int
	width  int
	height int
}

func VxNewModel(ch chan capture.VxConnection) VxModel {
	return VxModel{
		conns:  make([]capture.VxConnection, 0, 64),
		ch:     ch,
		counts: map[string]int{"SSH": 0, "FTP": 0, "HTTP": 0},
	}
}

func _wait(ch chan capture.VxConnection) tea.Cmd {
	return func() tea.Msg { return _connMsg(<-ch) }
}

func (m VxModel) Init() tea.Cmd {
	go capture.VxStartListener("SSH", 2222, m.ch)
	go capture.VxStartListener("FTP", 2121, m.ch)
	go capture.VxStartListener("HTTP", 8080, m.ch)
	return _wait(m.ch)
}

func (m VxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		if s := msg.String(); s == "q" || s == "ctrl+c" {
			return m, tea.Quit
		}
	case _connMsg:
		c := capture.VxConnection(msg)
		m.conns = append(m.conns, c)
		m.counts[c.Service]++
		return m, _wait(m.ch)
	}
	return m, nil
}

func _badge(svc string) string {
	switch svc {
	case "SSH":
		return _ssh.Render("SSH")
	case "FTP":
		return _ftp.Render("FTP")
	case "HTTP":
		return _http.Render("HTTP")
	}
	return svc
}

func (m VxModel) View() string {
	var b strings.Builder

	b.WriteString(_hdr.Render("☠  BEETRAP  ☠  | SSH·2222 | FTP·2121 | HTTP·8080") + "\n\n")

	total := m.counts["SSH"] + m.counts["FTP"] + m.counts["HTTP"]
	b.WriteString(fmt.Sprintf(
		"  %s %s   %s %s   %s %s   %s %s\n\n",
		_badge("SSH"), _bold.Render(fmt.Sprintf("%d", m.counts["SSH"])),
		_badge("FTP"), _bold.Render(fmt.Sprintf("%d", m.counts["FTP"])),
		_badge("HTTP"), _bold.Render(fmt.Sprintf("%d", m.counts["HTTP"])),
		_dim.Render("total"), _bold.Render(fmt.Sprintf("%d", total)),
	))

	maxRows := m.height - 10
	if maxRows < 1 {
		maxRows = 20
	}

	if len(m.conns) == 0 {
		b.WriteString(_dim.Render("  awaiting connections...") + "\n")
	} else {
		start := 0
		if len(m.conns) > maxRows {
			start = len(m.conns) - maxRows
		}
		for i := len(m.conns) - 1; i >= start; i-- {
			c := m.conns[i]
			ts := _dim.Render(c.Timestamp.Format("15:04:05"))
			payload := ""
			if c.Payload != "" {
				payload = _dim.Render("  » " + c.Payload)
			}
			b.WriteString(fmt.Sprintf("  %s  %s  %s%s\n", ts, _badge(c.Service), c.RemoteAddr, payload))
		}
	}

	b.WriteString("\n" + _dim.Render("  q quit"))
	return b.String()
}