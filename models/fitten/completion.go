package fitten

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gvcgo/fcode/cnf"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	completionURL = "https://fc.fittenlab.cn/codeapi/completion/generate_one_stage/"
)

func HandleAll(c *gin.Context) {
	mm, ok := c.Get(cnf.ModelCtxKey)
	if !ok {
		fmt.Println("no ai model found")
		return
	}

	model, ok := mm.(*cnf.AIModel)
	if !ok {
		fmt.Println("invalid ai model")
		return
	}

	apiKey := model.Key
	if apiKey == "" {
		apiKey = Login(model)
	}

	content, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	reqStr := string(content)
	body := bytes.NewBuffer(content)
	if strings.Contains(reqStr, cnf.DefaultConf.GetCursor()) {
		handleCompletions(c, body, apiKey)
	} else {
		handleChat(c, body, apiKey)
	}
}

func handleCompletions(c *gin.Context, body *bytes.Buffer, apiKey string) {
	req := &LspAIReq{}
	err := json.NewDecoder(body).Decode(req)
	if err != nil {
		return
	}

	if len(req.Messages) < 2 {
		return
	}

	var msg string
	cursor := cnf.DefaultConf.GetCursor()
	if cursor == "" {
		fmt.Println("no cursor specified!")
		os.Exit(1)
	}

	for _, m := range req.Messages {
		if m.Role == cnf.RoleUser {
			if strings.Contains(m.Content, cursor) {
				msg = m.Content
				break
			}
		}
	}

	if msg == "" {
		return
	}
	result := strings.SplitN(msg, cursor, 2)
	prefix := strings.TrimSpace(result[0])
	suffix := result[1]

	prompt := fmt.Sprintf(FCodeCompletionPrompt, prefix, suffix)
	payload := map[string]any{
		"inputs": prompt,
	}

	payloadJSON, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/%s?ide=%s&v=%s", completionURL, apiKey, IdeName, PluginVersion)
	resp, err := http.Post(url, "application/json", bytes.NewReader(payloadJSON))
	if err != nil {
		fmt.Println(err)
		return
	}

	r := &CompletionResponse{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	comp := strings.ReplaceAll(r.GeneratedText, "<.endoftext.>", "")
	comp = strings.ReplaceAll(comp, `<|endoftext|>`, "")

	lspAIResp := &cnf.CompResponse{
		Choices: []cnf.Choice{
			{
				FinishReason: "stop",
				Message: cnf.Message{
					Content: comp,
					Role:    cnf.RoleAssistant,
				},
			},
		},
	}

	c.JSON(http.StatusOK, lspAIResp)

}
