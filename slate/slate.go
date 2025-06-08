package slate

import (
	"maps"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keywordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204")).
			Background(lipgloss.Color("235"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type Model struct {
	activeTasks map[int]Model
	input       string
	idCounter   int
	textInput   textinput.Model
	finished    []string
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	for id, task := range m.activeTasks {
		m.activeTasks[id], cmd = task.Update(msg)
		switch cmd().(type) {
		case tea.QuitMsg:
			delete(m.activeTasks, id)
		default:
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var result string
	for task := range maps.Values(m.activeTasks) {
		result = lipgloss.JoinHorizontal(lipgloss.Left, result, task.View())
	}
	return keywordStyle.Render(result)
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}
