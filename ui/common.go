package ui

import (
	"github.com/muesli/termenv"
)

var (
	term = termenv.EnvColorProfile()

	subtle = makeFgStyle("241")
	dotStr = colorFg(" • ", "236")
)

type SpinningStatus int

const (
	Spinning SpinningStatus = iota
	Crashed
	UserCancel
	Finished
)

type JobCompletionCmd struct{}

func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}
