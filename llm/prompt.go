package llm

import (
	"log"
	"strings"

	"github.com/blackironj/easycommit/utils/restcli"
)

const (
	_defaultSummaryPrompt = `You are an expert developer specialist in creating commits.
Provide a super concise one sentence overall changes summary of the following "git diff" output following strictly the next rules:
- Do not use any code snippets, imports, file routes or bullets points.
- Do not mention the route of file that has been changes.
- Do not mention anything related to dependencies.
- Simply describe the MAIN GOAL of the changes.
- Output directly the summary in plain text.`

	_defaultCommitMsgPrompt = `Your only goal is to retrieve a single commit message.
Based on the following changes, combine them in ONE SINGLE commit message retrieving the global idea, following strictly the next rules:
- Always use the next format: "{type}: {commit_message}" where "{type}" is one of "feat", "fix", "docs", "style", "refactor", "test", "chore", "revert".
- Do not mention the route of files, name of files or imports path.
- Output directly only one commit message in plain text.
- Be as concise as possible. 72 characters max.
- Do not add only issues numeration nor explain your output.`
)

type Prompter struct {
	Endpoint    string
	Model       LLMmodel
	Temperature float64
	NumPredict  int
}

func NewPrompter(endpoint string, model LLMmodel, temperature float64, numPredict int) Prompter {
	return Prompter{
		Endpoint:    endpoint,
		Model:       model,
		Temperature: temperature,
		NumPredict:  numPredict,
	}
}

func (p Prompter) getTokenFromOllama(prompt string) string {
	cli := restcli.NewClient(restcli.Options{
		BaseUrl: p.Endpoint,
	})

	BodyParam := OllamaParams{
		Model:  p.Model,
		Prompt: prompt,
		Stream: false,
		Options: OllamaOptions{
			NumPredict:  p.NumPredict,
			Temperature: p.Temperature,
		},
	}

	var token OllamaToken
	_, err := cli.Post("/api/generate", &restcli.Params{Body: BodyParam}, &token)
	if err != nil {
		log.Print(err)
		return ""
	}
	return strings.TrimSpace(strings.Trim(token.Response, "\""))
}

func (p Prompter) GetSummary(diff, customSummaryPrompt string) string {
	prompt := customSummaryPrompt
	if customSummaryPrompt == "" {
		prompt = _defaultSummaryPrompt
	}

	prompt = prompt + "\n" + `Here is the "git diff" output:` + diff
	summary := p.getTokenFromOllama(prompt)

	return summary
}

func (p Prompter) GetCommitMsg(summary, customCommitPrompt string) string {
	prompt := customCommitPrompt
	if customCommitPrompt == "" {
		prompt = _defaultCommitMsgPrompt
	}

	prompt = prompt + "\n" + `Here are the summaries changes:` + summary
	commitMsg := p.getTokenFromOllama(prompt)

	return commitMsg
}
