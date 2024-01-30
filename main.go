package main

import (
	"fmt"
	"os/exec"

	"github.com/blackironj/easycommit/llm"
)

func main() {
	diff, err := exec.Command("git", "diff", "--cached", "--diff-algorithm=minimal").Output()
	if err != nil {
		panic(err)
	}

	prompter := llm.NewPrompter("http://127.0.0.1:11434", llm.Mistral, 0.7, 100)
	summary := prompter.GetSummary(string(diff), "")

	commit := prompter.GetCommitMsg(summary, "")
	fmt.Println("Generated commit msg : ", commit)

	if err := exec.Command("git", "commit", "-m", commit).Run(); err != nil {
		panic(err)
	}
}
