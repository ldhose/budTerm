package examples

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}
	spinnerHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

type SpinnerModel struct {
	index     int
	spinner   spinner.Model
	altScreen bool
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "h", "left":
			m.index--
			if m.index < 0 {
				m.index = len(spinners) - 1
			}
			m.resetSpinner()
			return m, m.spinner.Tick
		case "l", "right":
			m.index++
			if m.index >= len(spinners) {
				m.index = 0
			}
			m.resetSpinner()
			return m, m.spinner.Tick
		case " ":
			var cmd tea.Cmd
			if m.altScreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.altScreen = !m.altScreen
			return m, cmd

		default:
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *SpinnerModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}

func (m SpinnerModel) View() (str string) {
	var gap string
	switch m.index {
	case 1:
		gap = ""
	default:
		gap = " "
	}
	str = fmt.Sprintf("\n %s%s%s\n\n",
		m.spinner.View(), gap, textStyle("Spinning"))
	str += spinnerHelpStyle("h/l, ←/→: change spinner • q: exit\n")
	return
}

func StartSpinner() {
	m := SpinnerModel{}
	m.resetSpinner()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Could not run program", err)
		os.Exit(1)
	}
}
