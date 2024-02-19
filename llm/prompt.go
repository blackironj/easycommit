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

	_defaultCommitMsgPrompt = `You must recommend a list of at least three commit messages to the user.
Based on the following changes, combine them in ONE SINGLE commit message retrieving the global idea, following strictly the next rules:
- Always use the next format: "{type}: {commit_message}" where "{type}" is one of "feat", "fix", "docs", "style", "refactor", "test", "chore", "revert".
- Do not mention the route of files, name of files or imports path.
- Output directly commit messages line by line in plain text.
- No need numbers, bullet points, quote and double quote.
- Be as concise as possible. 72 characters max.
- Do not add only issues numeration nor explain your output.`
)

type Prompter struct {
	Endpoint string
	Model    LLMmodel
	Opt      OllamaOptions
}

func NewPrompter(endpoint string, model LLMmodel, opt ...OllamaOptions) Prompter {
	prompter := Prompter{
		Endpoint: endpoint,
		Model:    model,
	}
	if len(opt) > 0 {
		prompter.Opt = opt[0]
	}
	return prompter
}

func (p Prompter) getTokenFromOllama(prompt string) string {
	cli := restcli.NewClient(restcli.Options{
		BaseUrl: p.Endpoint,
	})

	BodyParam := OllamaParams{
		Model:   p.Model,
		Prompt:  prompt,
		Stream:  false,
		Options: p.Opt,
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
