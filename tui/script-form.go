package tui

import (
	"fmt"

	"scriptui/script"
	"scriptui/tui/constants"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ScriptFormModel struct {
	Script       script.Script
	TextInputs   []textinput.Model
	ActiveField  int
	err          error
	CallingModel tea.Model
}

const (
	width   = 40
	MaxLen  = 256
	hotPink = lipgloss.Color("#FF06B7")
	red     = lipgloss.Color("#FF0000")
)

var (
	borderStyle    = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(hotPink).Padding(0, 0)
	inputStyle     = lipgloss.NewStyle().Foreground(hotPink).Margin(2, 0)
	mandatoryStyle = lipgloss.NewStyle().Foreground(red).Margin(2, 0)
)

func NewScriptFormModel(script script.Script, callingModel tea.Model) *ScriptFormModel {
	var inputs []textinput.Model = make([]textinput.Model, len(script.Args))
	for i := range script.Args {
		inputs[i] = textinput.New()
		inputs[i].CharLimit = MaxLen
		inputs[i].Prompt = ""
		inputs[i].Width = width
		inputs[i].SetValue(script.Args[i].DefaultValue)
		inputs[i].Placeholder = script.Args[i].Placeholder
	}
	inputs[0].Focus()

	return &ScriptFormModel{Script: script, CallingModel: callingModel, TextInputs: inputs}
}

func (m *ScriptFormModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Scriptui :: Args")
}

func (m *ScriptFormModel) nextInput() {
	m.ActiveField++
	m.ActiveField %= len(m.Script.Args)
}

func (m *ScriptFormModel) prevInput() {
	m.ActiveField--
	if m.ActiveField < 0 {
		m.ActiveField = len(m.Script.Args) - 1
	}
}

func (m *ScriptFormModel) allMandatoryFilled() bool {
	for i, input := range m.TextInputs {
		if input.Value() == "" && m.Script.Args[i].Mandatory {
			return false
		}
	}
	return true
}

func (m *ScriptFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd = make([]tea.Cmd, len(m.TextInputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.allMandatoryFilled() {
				m.RunScript()
			}
		case tea.KeyEsc:
			return m.CallingModel, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN, tea.KeyDown:
			m.nextInput()
		}
		for i := range m.TextInputs {
			m.TextInputs[i].Blur()
		}
		m.TextInputs[m.ActiveField].Focus()
	}

	for i := range m.TextInputs {
		m.TextInputs[i], commands[i] = m.TextInputs[i].Update(msg)
	}

	commands = append(commands, tea.SetWindowTitle("klee++ :: Scripts :: Args"))
	return m, tea.Batch(commands...)
}

func (m *ScriptFormModel) RunScript() {
	for i := range m.Script.Args {
		if m.TextInputs[i].Value() == "false" || (!m.Script.Args[i].Mandatory && m.TextInputs[i].Value() == "") {
			m.Script.Args[i].Flag = ""
			m.Script.Args[i].Value = ""
		} else if m.TextInputs[i].Value() == "true" {
			continue
		} else {
			m.Script.Args[i].Value = m.TextInputs[i].Value()
		}
	}

	m.Script.RunScript()
}

func (m *ScriptFormModel) View() string {
	output := "" //line + "\n"

	for i, arg := range m.Script.Args {
		mandatory := ""
		if arg.Mandatory {
			mandatory += " * "
		}
		output += fmt.Sprintf("\t%s\n\t", inputStyle.Width(0).Render(arg.Title+mandatory))
		output += m.TextInputs[i].View()
		output += "\n"
	}

	output += m.helpView()
	return borderStyle.Render(output)
}

func (m *ScriptFormModel) helpView() string {
	return constants.HelpStyle("\n  ↑/shift+tab up • ↓/tab down • esc back • ctrl+c quit • enter run script (* mandatory field)\n ")
}
