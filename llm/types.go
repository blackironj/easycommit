package llm

type LLMmodel string

const (
	Mistral LLMmodel = "mistral"
)

type OllamaOptions struct {
	NumPredict  int     `json:"num_predict"`
	Temperature float64 `json:"temperature"`
}

type OllamaParams struct {
	Model   LLMmodel      `json:"model"`
	Prompt  string        `json:"prompt"`
	Stream  bool          `json:"stream"`
	Options OllamaOptions `json:"options"`
}

type OllamaToken struct {
	Model    LLMmodel `json:"model"`
	Response string   `json:"response"`
	Done     bool     `json:"done"`
}
