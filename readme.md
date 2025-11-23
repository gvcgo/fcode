## 什么是fcode？
- fcode是为lsp-ai项目做的一个adapter，支持fitten code，以及其他兼容OpenAI接口的大模型API。例如，美团的LongCat，阿里的ModelScope(qwen)，OpenAI等等。

## fcode带来什么好处？
- 原本不支持lsp-ai的fitten code可以在lsp-ai中使用。
- 在不同的Model之间快速无缝切换。

## 使用方法
- 编辑配置文件
```bash
mkdir -p ~/.fcode
cd ~/.fcode
touch conf.toml
# 复制conf_example/conf.toml 内容到上述文件，并修改username, password, key等等。可以根据自己的情况增减大模型配置。
```

## 配置lsp-ai
- 以helix editor为例
```text
# 见conf_example/helix_languages.toml
```

## 安装
```bash
go install github.com/moqsien/fcode
```

## 开启本地服务
```bash
fcode serve
```
