[中文](https://github.com/moqsien/fcode/blob/main/readme_cn.md)|[En](https://github.com/moqsien/fcode)

## What is fcode?
- fcode is an adapter for the [lsp‑ai](https://github.com/SilasMarvin/lsp-ai) project. It supports Fitten Code, Cloudflare AI Workers, and other large‑model APIs compatible with the OpenAI interface, such as Meituan’s LongCat, Alibaba’s ModelScope (qwen), Google’s Gemini, OpenAI, etc.

## What benefits does fcode bring?
- Fitten Code, which originally does not work with lsp‑ai, can now be used within lsp‑ai.
- Cloudflare AI Workers, previously unsupported, can also be used with lsp‑ai.
- Quickly and seamlessly switch between different models/APIs.
- Allows setting a local proxy for individual models (friendly for users in mainland China).

## How to use
- Edit the configuration file  

```bash
mkdir -p ~/.fcode
cd ~/.fcode
touch conf.toml
# Copy the contents of conf_example/conf.toml into this file and modify username, password, key, etc. You can add or remove model configurations as needed.
```

## Commands provided by fcode
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

## lsp‑ai conf
- Example for the [Helix editor](https://github.com/helix-editor/helix)  

```text
# See conf_example/helix_languages.toml
```

## Installation
```bash
go install github.com/moqsien/fcode@latest
```

## Start the local service
```bash
fcode stop && fcode serve>/dev/null 2>&1 &
```

## Gallery
- Fitten Code  

  ![fitten](https://github.com/moqsien/fcode/blob/main/imgs/lsp-ai_fitten.png)

- Chat  

  ![chat](https://github.com/moqsien/fcode/blob/main/imgs/lsp-ai_chat.png)
