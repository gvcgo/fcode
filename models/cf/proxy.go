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

	"github.com/gvcgo/fcode/cnf"

	"github.com/gin-gonic/gin"
)

type CFWorkersReq struct {
	Messages any `json:"messages"`
}

type Resust struct {
	Response string `json:"response"`
}

type CFWorkersResp struct {
	Result  Resust `json:"result"`
	Success bool   `json:"success"`
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

const (
	ReverseProxyErrCtxKey = "proxy_err_key"
)

func HandleAll(c *gin.Context) {
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

		reqBody := map[string]any{}
		_ = json.Unmarshal(bodyBuffer, &reqBody)
		messages := reqBody["messages"]

		bodyBuffer, _ = json.Marshal(&CFWorkersReq{
			Messages: messages,
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
		result := &CFWorkersResp{}
		err := json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}
		resp.Body.Close()

		lspAIResp := &cnf.CompResponse{
			Choices: []cnf.Choice{
				{
					FinishReason: "stop",
					Message: cnf.Message{
						Content: result.Result.Response,
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
