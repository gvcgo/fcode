package cf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/moqsien/fcode/cnf"
)

type CFOssReq struct {
	Model        string `json:"model"`
	Instructions string `json:"instructions"`
	Input        any    `json:"input"`
}

type CFOssContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type CFOssMsg struct {
	Type     string         `json:"type"`
	Role     string         `json:"role"`
	Status   string         `json:"status"`
	Contents []CFOssContent `json:"content"`
}

type CFOssResp struct {
	CFOssMsgs []CFOssMsg `json:"output"`
}

func HandleCFgptOss(c *gin.Context) {
	mm, ok := c.Get(cnf.ModelCtxKey)
	if !ok {
		fmt.Println("no model found")
		return
	}
	model, ok := mm.(*cnf.AIModel)
	if !ok {
		fmt.Println("invalid ai model")
		return
	}
	aiEndpoint, err := url.Parse(model.Api)
	if err != nil {
		fmt.Println("invalid api url: ", model.Api)
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(aiEndpoint)
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(r *http.Request) {
		// request body can only be read once, so we need to save it in a buffer.
		var bodyBuffer []byte
		if r.Body != nil {
			var err error
			bodyBuffer, err = io.ReadAll(r.Body)
			if err != nil {
				ctx := context.WithValue(r.Context(), ReverseProxyErrCtxKey, err)
				*r = *r.WithContext(ctx)
				return
			}
		}

		reqBody := map[string][]cnf.Message{}
		_ = json.Unmarshal(bodyBuffer, &reqBody)

		messages := []cnf.Message{}
		msgs := reqBody["messages"]
		instructions := ""
		for _, m := range msgs {
			if m.Role == cnf.RoleSystem {
				instructions = m.Content
			} else {
				messages = append(messages, m)
			}
		}
		bodyBuffer, _ = json.Marshal(&CFOssReq{
			Model:        model.Model,
			Instructions: instructions,
			Input:        messages,
		})

		originalDirector(r)
		if len(bodyBuffer) > 0 {
			r.Body = io.NopCloser(bytes.NewReader(bodyBuffer))
			r.GetBody = func() (io.ReadCloser, error) {
				return r.Body, nil
			}
			r.ContentLength = int64(len(bodyBuffer))
		}

		r.Host = aiEndpoint.Host
		r.URL.Path = aiEndpoint.Path
		r.Header.Del("Origin")
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", model.Key))
	}

	reverseProxy.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if proxyErr, ok := req.Context().Value(ReverseProxyErrCtxKey).(error); ok {
			return nil, proxyErr
		}
		// 设置本地代理
		localProxy := c.GetString(cnf.ProxyCtxKey)
		if localProxy != "" && model.UseProxy {
			proxyURL, _ := url.Parse(localProxy)
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			return transport.RoundTrip(req)
		}

		return http.DefaultTransport.RoundTrip(req)
	})

	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) error {
		result := &CFOssResp{}
		err := json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}
		resp.Body.Close()

		var res *CFOssContent
	OUTTER:
		for _, o := range result.CFOssMsgs {
			if o.Role == cnf.RoleAssistant {
				for _, cc := range o.Contents {
					if cc.Text != "" && cc.Type == "output_text" {
						res = &cc
						break OUTTER
					}
				}
			}
		}
		lspAIResp := &cnf.CompResponse{
			Choices: []cnf.Choice{
				{
					FinishReason: "stop",
					Message: cnf.Message{
						Content: res.Text,
						Role:    cnf.RoleAssistant,
					},
				},
			},
		}

		body, err := json.Marshal(lspAIResp)
		if err != nil {
			return err
		}

		resp.Header.Set("X-Proxy-Processed", "true")
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		resp.ContentLength = int64(len(body))
		return nil
	}

	reverseProxy.ServeHTTP(c.Writer, c.Request)
}
