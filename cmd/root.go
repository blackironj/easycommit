package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/blackironj/easycommit/llm"
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
		diff, err := exec.Command("git", "diff", "--cached", "--diff-algorithm=minimal").Output()
		if err != nil {
			panic(err)
		}

		prompter := llm.NewPrompter(endpoint, llm.LLMmodel(llmModel), llmTemperature, llmNumPredict)
		summary := prompter.GetSummary(string(diff), "")

		commit := prompter.GetCommitMsg(summary, "")

		if err := exec.Command("git", "commit", "-m", commit).Run(); err != nil {
			panic(err)
		}
	},
}

func SetGlobalFlag() {
	RootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", llm.DefaultOllamaEndpoint, "ollama host url (default: localhost)")
	RootCmd.PersistentFlags().StringVarP(&llmModel, "model", "m", string(llm.Mistral), "llama model (default: mistral)")
	RootCmd.PersistentFlags().Float64VarP(&llmTemperature, "temperature", "t", llm.DefaultTemperature, "temperature (default: 0.7)")
	RootCmd.PersistentFlags().IntVarP(&llmNumPredict, "num-predict", "n", llm.DefaultNumPredict, "num predict (default: 100)")
}
