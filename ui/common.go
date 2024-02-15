package ui

type SpinningStatus int

const (
	Spinning SpinningStatus = iota
	Crashed
	UserCancel
	Finished
)

type JobCompletionCmd struct{}
