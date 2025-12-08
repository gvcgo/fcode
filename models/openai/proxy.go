package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/moqsien/fcode/cnf"

	"github.com/gin-gonic/gin"
)

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

	proxy := httputil.NewSingleHostReverseProxy(aiEndpoint)
	originalDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		fmt.Println(r.Header)
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
		if _, ok := reqBody["model"]; ok {
			reqBody["model"] = model.Model
		}

		if strings.Contains(aiEndpoint.Host, "google") {
			delete(reqBody, "frequency_penalty")
			// delete(reqBody, "presence_penalty")
		}

		bodyBuffer, _ = json.Marshal(reqBody)
		// fmt.Println("------> ", string(bodyBuffer))

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

	proxy.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if proxyErr, ok := req.Context().Value(ReverseProxyErrCtxKey).(error); ok {
			return nil, proxyErr
		}
		return http.DefaultTransport.RoundTrip(req)
	})

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// fmt.Println(resp.StatusCode)
		// fmt.Println("==> ", string(body))
		resp.Body.Close()
		resp.Header.Set("X-Proxy-Processed", "true")
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		resp.ContentLength = int64(len(body))
		return nil
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
