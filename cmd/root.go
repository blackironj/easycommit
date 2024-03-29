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
	llmNumPredict  int64
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

		if string(diff) == "" {
			fmt.Println("😔 no diff")
			return
		}

		var summary string
		var prompter llm.Prompter
		var resultCommits []string

		spinningProg := ui.NewSpinningProgram()
		go func() {
			prompter = llm.NewPrompter(endpoint,
				llm.LLMmodel(llmModel),
				llm.OllamaOptions{
					Temperature: &llmTemperature,
					NumPredict:  &llmNumPredict,
				})
			summary = prompter.GetSummary(string(diff), "")
			rawCommits := prompter.GetCommitMsg(summary, "")

			rawCommits = strings.Replace(rawCommits, "\n\n", "\n", -1)
			resultCommits = strings.Split(rawCommits, "\n")
			for i := range resultCommits {
				resultCommits[i] = strings.Trim(resultCommits[i], "\"")
			}
			spinningProg.Send(ui.JobCompletionCmd{})
		}()

		switch ui.RunSpinning(spinningProg) {
		case ui.UserCancel:
			fmt.Println("❌ user cancel")
			return
		case ui.Crashed:
			fmt.Println("😔 crashed")
			return
		}

		interactionProg := ui.NewMsgChoiceProgram(resultCommits)
		selectedCommit := ui.RunMsgChoice(interactionProg)

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
	RootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", llm.DefaultOllamaEndpoint, "ollama host url")
	RootCmd.PersistentFlags().StringVarP(&llmModel, "model", "m", string(llm.Mistral), "llama model")
	RootCmd.PersistentFlags().Float64VarP(&llmTemperature, "temperature", "t", llm.DefaultTemperature, "temperature")
	RootCmd.PersistentFlags().Int64VarP(&llmNumPredict, "num-predict", "n", llm.DefaultNumPredict, "num predict")

	RootCmd.PersistentFlags().BoolVarP(&logger.DebugFlag, "debug", "d", false, "for debugging log")
}
