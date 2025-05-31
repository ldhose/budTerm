package examples

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type ScreenModel struct {
	altScreen  bool
	quitting   bool
	suspending bool
}

func (m ScreenModel) Init() tea.Cmd {
	return nil
}

func (m ScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.ResumeMsg:
		m.suspending = false
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+z":
			m.suspending = true
			return m, tea.Suspend
		case " ":
			var cmd tea.Cmd
			if m.altScreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.altScreen = !m.altScreen
			return m, cmd
		}
	}
	return m, nil
}

func StartScreen() {
	if _, err := tea.NewProgram(ScreenModel{}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m ScreenModel) View() string {
	if m.suspending {
		return ""
	}

	if m.quitting {
		return "Bye! \n"
	}

	const (
		altScreenMode = " altscreen mode "
		inlineMode    = " inlineMode "
	)

	var mode string
	if m.altScreen {
		mode = altScreenMode
	} else {
		mode = inlineMode
	}
	return fmt.Sprintf("\n\n You are in %s \n\n\n", keywordStyle.Render(mode)) +
		helpStyle.Render(" space: switch modes • ctrl-z: suspend • q: exit\n")

}
