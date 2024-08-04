package constants

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var WindowSize tea.WindowSizeMsg

var DocStyle = lipgloss.NewStyle().Margin(0, 2)

var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd534b")).Render

type keymap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
