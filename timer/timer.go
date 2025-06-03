package timer

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keywordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204")).
			Background(lipgloss.Color("235"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

const (
	Focus = iota
	ShortBreak
	LongBreak
	Paused
	Finished
	Running
)

type ViewType uint8

const (
	tagView ViewType = iota
	timerView
)
const (
	altScreenMode = " altscreen mode "
	inlineMode    = " inlineMode "
)

var stateName = map[int]string{
	Focus:      "Focus",
	ShortBreak: "Short Break",
	LongBreak:  "Long Break",
	Paused:     "Paused",
}

type Model struct {
	//All time are in seconds
	ID       uint16
	tag      uint16
	counter  uint16
	name     string
	state    int
	duration uint16
}

func (m *Model) Reset() {
	m.counter = m.duration
	m.state = Running
}

func NewTimer(timeout uint16, id uint16) Model {
	return Model{
		ID:       id,
		tag:      id,
		counter:  timeout,
		state:    Running,
		duration: timeout,
	}
}

func StartTimer(timerModel tea.Model) {
	if _, err := tea.NewProgram(timerModel, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			if m.state == Paused {
				m.state = Running
				break
			}
			if m.state == Running {
				m.state = Paused
				break
			}
			return m, tea.Batch(m.tick())
		}
	case TickMsg:
		if msg.ID != 0 && msg.ID != m.ID {
			return m, nil
		}
		if msg.Finished {
			break
		}
		if m.state == Paused {
			break
		}
		if m.counter <= 0 {
			break
		}
		m.counter -= 1
		return m, tea.Batch(m.tick())
	}
	return m, nil
}

func (m Model) View() string {
	if m.state == Finished {
		return keywordStyle.Render(("__ : __"))
	}
	min := m.counter / 60
	sec := m.counter % 60
	result := fmt.Sprintf("%02d : %02d", min, sec)
	return keywordStyle.Render(result)
}

type TimeoutMsg struct {
	ID uint16
}

func (m Model) Finished() bool {
	return m.counter <= 0
}

func (m *Model) Toggle() {
	if m.state == Paused {
		m.state = Running
	}
	if m.state == Running {
		m.state = Paused
	}
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return TickMsg{ID: m.ID, tag: m.tag, Finished: m.Finished()}
	})
}

type TickMsg struct {
	ID       uint16
	tag      uint16
	Finished bool
}
