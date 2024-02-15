package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/blackironj/easycommit/utils/logger"
)

type spinningModel struct {
	spinner spinner.Model
	status  SpinningStatus
}

func NewSpinningProgram() *tea.Program {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return tea.NewProgram(spinningModel{
		spinner: s,
		status:  Spinning,
	})
}

func RunSpinning(prog *tea.Program) SpinningStatus {
	resModel, err := prog.Run()
	if err != nil {
		logger.Get().Err(err)
		return Crashed
	}
	return resModel.(spinningModel).status
}

func (m spinningModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinningModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.status = UserCancel
			return m, tea.Quit
		default:
			return m, nil
		}

	case JobCompletionCmd:
		m.status = Finished
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinningModel) View() string {
	str := fmt.Sprintf("\n%s Generating...\n\n", m.spinner.View())
	if m.status == Finished {
		return str + "\n✅ Job finished! \n\n"
	}
	return str
}
