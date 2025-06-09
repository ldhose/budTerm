package task

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	store "github.com/ldhose/budTerm/task/storage"
	"github.com/ldhose/budTerm/timer"
)

type TaskModel struct {
	timerModel timer.Model
	tagsModel  textinput.Model
	fullScreen bool
	name       string
}

type instructionType uint8

const (
	Timer instructionType = iota
)

type instruction struct {
	instructionType instructionType
	name            string
	value           string
}

func (m TaskModel) Init() tea.Cmd {
	return tea.Sequence(m.tagsModel.Focus(), m.timerModel.Init())
}

func (m TaskModel) Update(msg tea.Msg) (TaskModel, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.tagsModel, cmd = m.tagsModel.Update(msg)
	cmds = append(cmds, cmd)
	//TODO store results when timer is finished.
	m.timerModel, cmd = m.timerModel.Update(msg)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case timer.TimeoutMsg:
		// log.Println("Timer has finished running")
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			m.timerModel.Toggle()
			m.timerModel, cmd = m.timerModel.Update(m.timerModel.TickMsg())
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case "tab":
			m.tagsModel.Focus()
		case "ctrl+q":
			return m, tea.Quit
		case "ctrl+l": // rest
			return m, tea.Batch(cmds...)
		case "ctrl+b": // break
			return m, tea.Batch(cmds...)
		case "ctrl+r":
			m.timerModel.Reset()
			return m, tea.Batch(cmds...)
		case "enter":
			input := m.tagsModel.Value()
			processInput(input, &m)
			m.name = input
			m.tagsModel.Reset()
		}
	}

	return m, tea.Batch(cmds...)
}

func processInput(input string, m *TaskModel) {
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

func execute(command instruction, m *TaskModel) {
	//TODO store info to file after processing.
	_, err := strconv.ParseUint(command.value, 10, 8)
	if err == nil {
		if command.name != "" {
			m.tagsModel.Reset()
			m.tagsModel.Blur()
			// m.timerModel = timer.NewTimer(uint16(value), 1, command.name)
			// newTask := TaskModel{
			// 	timerModel: m.timerModel,
			// 	tagsModel:  m.tagsModel,
			// 	name:       command.name,
			// }
			// StartTask(newTask)
			store.Append(command.name)
		}
	}
}

func Nothing(result string) {

}

func (m TaskModel) Finished() bool {
	return m.timerModel.Finished()
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

// func StartApp() {
// 	// newTimer := timer.NewTimer(20, 1, msg)
// 	newTag := textinput.New()

// 	if _, err := tea.
// 		NewProgram(
// 			TaskModel{
// 				// name:       msg,
// 				// timerModel: newTimer,
// 				tagsModel: newTag},
// 			tea.WithAltScreen()).
// 		Run(); err != nil {
// 		fmt.Println("Error running program:", err)
// 		os.Exit(1)
// 	}
// }

// func StartTask(task TaskModel) {
// 	if _, err := tea.
// 		NewProgram(
// 			task,
// 			tea.WithAltScreen()).
// 		Run(); err != nil {
// 		fmt.Println("Error running program:", err)
// 		os.Exit(1)
// 	}

// }
