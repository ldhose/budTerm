package task

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ldhose/budTerm/timer"
)

type TaskModel struct {
	timerModel timer.Model
	tagsModel  textinput.Model
	fullScreen bool
}

func (m TaskModel) Init() tea.Cmd {
	return tea.Sequence(m.tagsModel.Focus(), m.timerModel.Init())
}

func (m TaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.tagsModel, cmd = m.tagsModel.Update(msg)
	cmds = append(cmds, cmd)
	m.timerModel, cmd = m.timerModel.Update(msg)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.tagsModel.Focus()
		case "q":
			return m, tea.Quit
		case " ":
			var cmd tea.Cmd
			if m.fullScreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.fullScreen = !m.fullScreen
			return m, cmd
		case "enter":
			m.tagsModel.Reset()
			m.tagsModel.Blur()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m TaskModel) View() string {
	var s string
	s += lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("%50s", m.timerModel.View()),
		fmt.Sprintf("%4s", m.tagsModel.View()))
	return s
}

func (m TaskModel) Trap() string {
	return "trap"
}

func StartTask() {
	newTimer := timer.NewTimer(20, 1)
	newTag := textinput.New()

	if _, err := tea.
		NewProgram(
			TaskModel{
				timerModel: newTimer,
				tagsModel:  newTag},
			tea.WithAltScreen()).
		Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
