package fitten

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gvcgo/fcode/cnf"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	loginURL     = "https://fc.fittenlab.cn/codeuser/login"
	ficoTokenURL = "https://fc.fittenlab.cn/codeuser/get_ft_token"
)

func Login(model *cnf.AIModel) string {
	client := &http.Client{Timeout: 10 * time.Second}

	// Step 1: Login to get user token
	loginData := map[string]string{"username": model.Username, "password": model.Password}
	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var loginResp struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.NewDecoder(bytes.NewBuffer(content)).Decode(&loginResp)

	if loginResp.Code != 200 {
		fmt.Println(string(content))
		os.Exit(1)
	}

	userToken := loginResp.Data.Token

	// Step 2: Get FICO token (API key)
	req2, _ := http.NewRequest("GET", ficoTokenURL, nil)
	req2.Header.Set("Authorization", "Bearer "+userToken)
	resp2, err := client.Do(req2)
	if err != nil || resp2.StatusCode != http.StatusOK {
		fmt.Println("Get token failed")
		os.Exit(1)
	}
	defer resp2.Body.Close()

	var ficoResp struct {
		Data struct {
			FicoToken string `json:"fico_token"`
		} `json:"data"`
	}
	json.NewDecoder(resp2.Body).Decode(&ficoResp)
	cnf.DefaultConf.SetApiKey(model.Name, ficoResp.Data.FicoToken)

	return ficoResp.Data.FicoToken
}
