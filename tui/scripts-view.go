package tui

import (
	"fmt"
	"os"

	"scriptui/script"
	"scriptui/tui/constants"

	"github.com/charmbracelet/bubbles/key"
	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type scriptItem struct {
	title, desc string
}

func (i scriptItem) Title() string       { return i.title }
func (i scriptItem) Description() string { return i.desc }
func (i scriptItem) FilterValue() string { return i.title }

type ScriptsView struct {
	ScriptMap   map[string]*script.Script
	list        list.Model
	parentModel tea.Model
}

func requestWindowSize() tea.Msg {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return nil
	}
	return tea.WindowSizeMsg{Width: width, Height: height}
}

func NewList(scriptMap map[string]*script.Script) list.Model {
	var scriptsList []list.Item
	for _, s := range scriptMap {
		scriptsList = append(scriptsList, scriptItem{title: s.Name, desc: s.Desc})
	}
	l := list.New(scriptsList, list.NewDefaultDelegate(), 30, 30)
	l.Title = "scriptui"
	return l
}

func (sv *ScriptsView) Init() tea.Cmd {
	scripts, err := script.LoadScripts()
	if err != nil {
		panic(err)
	}
	sv.ScriptMap = scripts
	sv.list = NewList(scripts)
	return tea.Batch(tea.SetWindowTitle("Scriptui"), requestWindowSize)
}

func (sv *ScriptsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			return sv, tea.Quit
		case key.Matches(msg, constants.Keymap.Enter):
			selectedItem := sv.list.SelectedItem()
			if selectedItem != nil {
				selectedItem := selectedItem.(scriptItem) // type assertion
				scriptX := sv.ScriptMap[selectedItem.title]
				if len(scriptX.Args) == 0 {
					// This call will end the program, and therefore never returns.
					scriptX.RunScript()
				}
				scriptModel := NewScriptFormModel(*scriptX, sv)
				return scriptModel, scriptModel.Init()
			}
		}

	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		sv.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	}

	var cmd tea.Cmd
	sv.list, cmd = sv.list.Update(msg)
	return sv, tea.Batch(cmd, tea.SetWindowTitle("klee++ :: Scripts"))
}

var docStyle = lipgloss.NewStyle().Margin(0, 2)

func (sv *ScriptsView) View() string {
	return docStyle.Render(sv.list.View())
}
