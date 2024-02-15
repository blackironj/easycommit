package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blackironj/easycommit/llm"
	"github.com/blackironj/easycommit/ui"
	"github.com/blackironj/easycommit/utils/logger"
)

var (
	endpoint       string
	llmModel       string
	llmTemperature float64
	llmNumPredict  int
)

var RootCmd = &cobra.Command{
	Use:   "easycommit",
	Short: "simple cli tool for generating commit using Ollama",
	Run: func(cmd *cobra.Command, args []string) {
		l := logger.Get()

		diff, err := exec.Command("git", "diff", "--cached", "--diff-algorithm=minimal").Output()
		if err != nil {
			fmt.Println("😔 git diff failed")
			l.Err(err).Send()
			return
		}

		var summary string
		var prompter llm.Prompter
		var resultCommits []string

		spinningProg := ui.NewSpinningProgram()
		go func() {
			prompter = llm.NewPrompter(endpoint, llm.LLMmodel(llmModel), llmTemperature, llmNumPredict)
			summary = prompter.GetSummary(string(diff), "")
			rawCommits := prompter.GetCommitMsg(summary, "")

			rawCommits = strings.Replace(rawCommits, "\n\n", "\n", -1)
			resultCommits = strings.Split(rawCommits, "\n")
			for i := range resultCommits {
				resultCommits[i] = strings.Trim(resultCommits[i], "\"")
			}
			spinningProg.Send(ui.JobCompletionCmd{})
		}()

		ui.RunSpinning(spinningProg)

		interactionProg := ui.NewInteractionProgram(resultCommits)
		selectedCommit := ui.RunInteraction(interactionProg)

		l.Debug().Strs("recommended", resultCommits).Str("selected", selectedCommit).Send()
		if selectedCommit == "" {
			fmt.Println("❌ user cancel")
			return
		}

		if err := exec.Command("git", "commit", "-m", selectedCommit).Run(); err != nil {
			fmt.Println("😔 commit failed")
			l.Err(err).Send()
			return
		}
		fmt.Println("😃 commit success!")
	},
}

func SetGlobalFlag() {
	RootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", llm.DefaultOllamaEndpoint, "ollama host url (default: localhost)")
	RootCmd.PersistentFlags().StringVarP(&llmModel, "model", "m", string(llm.Mistral), "llama model (default: mistral)")
	RootCmd.PersistentFlags().Float64VarP(&llmTemperature, "temperature", "t", llm.DefaultTemperature, "temperature (default: 0.7)")
	RootCmd.PersistentFlags().IntVarP(&llmNumPredict, "num-predict", "n", llm.DefaultNumPredict, "num predict (default: 100)")

	RootCmd.PersistentFlags().BoolVarP(&logger.DebugFlag, "debug", "d", false, "for debugging log (default: false)")
}
