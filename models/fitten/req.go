package fitten

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

type FCodeDelta struct {
	Delta string `json:"delta"`
}
