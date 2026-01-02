package cnf

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type CompResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Logprobs     *string `json:"logprobs"`
	Message      Message `json:"message"`
}

type Message struct {
	Content          string  `json:"content"`
	ReasoningContent string  `json:"reasoning_content"`
	ReasoningDetails *string `json:"reasoning_details"`
	Role             string  `json:"role"`
	TaskID           *string `json:"task_id"`
}
