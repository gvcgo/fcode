package main

import "fmt"

const (
	completionURL = "https://fc.fittenlab.cn/codeapi/completion/generate_one_stage/"
	apiKeyFile    = ".vimapikey"
	ideName       = "vim"
	pluginVersion = "0.2.1"
)

func main() {
	fmt.Println(DefaultKey.APIKey)
}
