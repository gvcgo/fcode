package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	completionURL = "https://fc.fittenlab.cn/codeapi/completion/generate_one_stage/"
)

func handleAll(c *gin.Context) {
	content, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	reqStr := string(content)
	body := bytes.NewBuffer(content)
	if strings.Contains(reqStr, DefaultConf.Cursor) {
		handleCompletions(c, body)
	} else {
		handleChat(c, body)
	}
}

func handleCompletions(c *gin.Context, body *bytes.Buffer) {
	req := &LspAIReq{}
	err := json.NewDecoder(body).Decode(req)
	if err != nil {
		return
	}

	if len(req.Messages) < 2 {
		return
	}

	var msg string
	cursor := DefaultConf.GetCursor()
	if cursor == "" {
		fmt.Println("no cursor specified!")
		os.Exit(1)
	}

	for _, m := range req.Messages {
		if m.Role == RoleUser {
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

	prompt := fmt.Sprintf("!FCPREFIX!%s!FCSUFFIX!%s!FCMIDDLE!", prefix, suffix)
	payload := map[string]any{
		"inputs": prompt,
	}

	payloadJSON, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/%s?ide=%s&v=%s", completionURL, DefaultKey.APIKey, IdeName, PluginVersion)
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

	lspAIResp := &Response{
		[]Choice{
			{
				FinishReason: "stop",
				Message: Message{
					Content: comp,
					Role:    RoleAssistant,
				},
			},
		},
	}

	c.JSON(http.StatusOK, lspAIResp)

}
