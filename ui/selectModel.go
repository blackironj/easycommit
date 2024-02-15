package ui

import (
	"fmt"

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
	return choicesView(m)
}

func choicesView(m interactionModel) string {
	selected := m.cursor
	tpl := "* Which commit would you like? *\n\n"
	tpl += "%s\n"
	tpl += subtle("up/down: select") + dotStr + subtle("enter: choose") + dotStr + subtle("ctrl+c, q, esc: quit")

	var choices string
	for i, elem := range m.choiceList {
		choices += checkbox(elem, selected == i) + "\n"
	}
	return fmt.Sprintf(tpl, choices)
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return "[ ] " + label
}
