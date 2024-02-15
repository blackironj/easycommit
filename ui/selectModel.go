package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/blackironj/easycommit/utils/logger"
)

type interactionModel struct {
	cursor     int
	selected   string
	choiceList []string
}

func NewInteractionProgram(choiceList []string) *tea.Program {
	return tea.NewProgram(interactionModel{
		cursor:     0,
		selected:   "",
		choiceList: choiceList,
	})
}

func RunInteraction(prog *tea.Program) string {
	resModel, err := prog.Run()
	if err != nil {
		logger.Get().Err(err)
		return ""
	}

	finishedModel, ok := resModel.(interactionModel)
	if ok && finishedModel.selected != "" {
		return finishedModel.selected
	}
	return ""
}

func (m interactionModel) Init() tea.Cmd {
	return nil
}

func (m interactionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.selected = m.choiceList[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choiceList) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choiceList) - 1
			}
		}
	}
	return m, nil
}

func (m interactionModel) View() string {
	s := strings.Builder{}
	s.WriteString("* Which commit would you like? *\n\n")

	for i := 0; i < len(m.choiceList); i++ {
		if m.cursor == i {
			s.WriteString("[x] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(m.choiceList[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
