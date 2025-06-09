package slate

import (
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ldhose/budTerm/task"
	"github.com/ldhose/budTerm/timer"
)

var (
	keywordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204")).
			Background(lipgloss.Color("235"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type instructionType uint8

const (
	Timer instructionType = iota
)

type instruction struct {
	instructionType instructionType
	name            string
	value           string
}

type Model struct {
	activeTasks map[int]task.TaskModel
	input       string
	idCounter   int
	textBox     textinput.Model
	finished    []string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	var cmds []tea.Cmd
	for t := range maps.Values(m.activeTasks) {
		cmds = append(cmds, t.Init())
	}
	return tea.Batch(cmds...)
}

func StartApp() {
	// newTimer := timer.NewTimer(20, 1, msg)
	// newTag := textinput.New()
	newTag := textinput.New()
	if _, err := tea.
		NewProgram(
			Model{textBox: newTag},
			tea.WithAltScreen()).
		Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func processInput(input string, m *Model) {
	parts := strings.Split(input, ",")
	if len(parts) == 1 {
		return
	}
	var command instruction
	if len(parts) == 2 {
		command = instruction{
			instructionType: Timer,
			name:            parts[1],
		}
	}

	if len(parts) == 3 {
		command = instruction{
			instructionType: Timer,
			name:            parts[1],
			value:           parts[2],
		}
	}
	execute(command, m)

}

func execute(command instruction, m *Model) {
	//TODO store info to file after processing.
	value, err := strconv.ParseUint(command.value, 10, 8)
	if err == nil {
		if command.name != "" {
			m.textBox.Reset()
			m.textBox.Blur()
			m.timerModel = timer.NewTimer(uint16(value), 1, command.name)
			newTask := TaskModel{
				timerModel: m.timerModel,
				tagsModel:  m.tagsModel,
				name:       command.name,
			}
			// StartTask(newTask)
			store.Append(command.name)
		}
	}
}
