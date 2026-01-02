package fitten

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gvcgo/fcode/cnf"
)

const (
	chatUrl      = "https://codewebchat.fittenlab.cn/codeapi/chat"
	MsgSystem    = "<|system|>\n你是一个分析师，请根据聊天内容做出周密的分析，并给出精确的回答。\n<|end|>"
	MsgUser      = "<|user|>\n%s\n<|end|>"
	MsgAssistant = "<|assistant|>\n%s\n<|end|>"
	MsgTail      = "\n<|assistant|>"
)

func prepareInput(msgs ...Msg) string {
	fcodeMsgs := []string{
		MsgSystem,
	}
	for _, m := range msgs {
		if m.Content == "" {
			continue
		}

		switch m.Role {
		case cnf.RoleUser:
			fcodeMsgs = append(fcodeMsgs, fmt.Sprintf(MsgUser, m.Content))
		case cnf.RoleAssistant:
			fcodeMsgs = append(fcodeMsgs, fmt.Sprintf(MsgAssistant, m.Content))
		}
	}

	fcodeMsgs = append(fcodeMsgs, MsgTail)
	return strings.Join(fcodeMsgs, "\n")
}

func handleChat(c *gin.Context, body *bytes.Buffer, apiKey string) {
	req := &LspAIReq{}
	err := json.NewDecoder(body).Decode(req)
	if err != nil {
		return
	}

	if len(req.Messages) == 0 {
		return
	}

	inputs := prepareInput(req.Messages...)

	payload := map[string]any{
		"ft_token": apiKey,
		"inputs":   inputs,
	}
	payloadJSON, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s?ide=%s&show_shortcut=0&apikey=%s", chatUrl, IdeName, apiKey)
	resp, err := http.Post(url, "application/json", bytes.NewReader(payloadJSON))
	if err != nil {
		fmt.Println(err)
		return
	}

	var result string
	ss := bufio.NewScanner(resp.Body)
	for ss.Scan() {
		line := ss.Bytes()
		if len(line) == 0 {
			continue
		}
		var d FCodeDelta
		if err := json.Unmarshal(line, &d); err != nil {
			fmt.Println(err)
			continue
		}
		result += d.Delta
	}

	lspAIResp := &cnf.CompResponse{
		Choices: []cnf.Choice{
			{
				FinishReason: "stop",
				Message: cnf.Message{
					Content: result,
					Role:    cnf.RoleAssistant,
				},
			},
		},
	}

	c.JSON(http.StatusOK, lspAIResp)
}
