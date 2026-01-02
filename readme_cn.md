## 什么是fcode？
- fcode是为[lsp-ai](https://github.com/SilasMarvin/lsp-ai)项目做的一个adapter，支持fitten code，cloudflare AI workers，以及其他兼容OpenAI接口的大模型API。例如，美团的LongCat，阿里的ModelScope(qwen)，google的gemini，OpenAI等等。

## fcode带来什么好处？
- 原本不支持lsp-ai的fitten code可以在lsp-ai中使用。
- 原本不支持的cloudflare AI workers也可以在lsp-ai中使用。
- 在不同的Model/API之间快速无缝切换。
- 支持为单个模型设置本地代理(对于国内用户友好)

## 使用方法
- 编辑配置文件
```bash
mkdir -p ~/.fcode
cd ~/.fcode
touch conf.toml
# 复制conf_example/conf.toml 内容到上述文件，并修改username, password, key等等。可以根据自己的情况增减大模型配置。
```

## fcode提供了哪些命令
```bash
fcode -h   
Usage:
   [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List available model names.
  serve       Run server for lsp-ai.
  show        Show config file path.
  stop        Stop fcode server.
  use         Use an available model.

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```

## 配置lsp-ai
- 以[helix editor](https://github.com/helix-editor/helix)为例
```text
# 见conf_example/helix_languages.toml
```

## 安装
```bash
go install github.com/moqsien/fcode@latest
```

## 开启本地服务
```bash
fcode stop && fcode serve>/dev/null 2>&1 &
```

## Gallery
- fitten code

  ![fitten](https://github.com/moqsien/fcode/blob/main/imgs/lsp-ai_fitten.png)

- chat
  ![chat](https://github.com/moqsien/fcode/blob/main/imgs/lsp-ai_chat.png)
