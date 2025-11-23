package fitten

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type Msg struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type LspAIReq struct {
	Messages []Msg  `json:"messages"`
	Model    string `json:"model"`
}

type CompletionResponse struct {
	GeneratedText string `json:"generated_text"`
}

type CompResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Logprobs     *string `json:"logprobs"`
	Message      Message `json:"message"`
}

// Message 包含助手回复的详细信息
type Message struct {
	Content          string  `json:"content"`
	ReasoningContent string  `json:"reasoning_content"`
	ReasoningDetails *string `json:"reasoning_details"`
	Role             string  `json:"role"`
	TaskID           *string `json:"task_id"`
}

type FCodeDelta struct {
	Delta string `json:"delta"`
}
