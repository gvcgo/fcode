alias ls="eza"
alias ga="git add"
alias gc="git commit"
alias gp="git push"
alias gpl="git pull"
alias gst="git status"
alias gsw="git switch -C"
alias gswt="git switch --track"
alias gbde="git branch -d"
alias gbl="git branch -l"
alias gch="git checkout"
alias gco="git checkout ."
alias gd="git diff"
alias gpmm="git fetch origin && git merge origin/master"

# If you come from bash you might have to change your $PATH.
export PATH=$HOME/bin:$HOME/.local/bin:/usr/local/bin:$PATH
export DOCKER_PATH="/Applications/Docker.app/Contents/Resources/bin"
export PATH="$PATH:$DOCKER_PATH"

docker() {
    if ! pgrep -x "Docker Desktop" > /dev/null; then
        echo "Docker Desktop not running..."
        open --background -a "Docker Desktop"
        local timeout=60
        local count=0
        while ! docker info > /dev/null 2>&1; do
            ((count++))
            if [ $count -ge $timeout ]; then
                echo "error: Docker Desktop starting timedout(waited for ${timeout}s))"
                return 1
            fi
            sleep 2
            echo "waiting for Docker Desktop starting... ($count/$timeout)"
        done
        echo "Docker Desktop is ready!"
    fi

    "$DOCKER_PATH/docker" "$@"
}

# z for zoxide
eval "$(zoxide init zsh)"

if [ -z "$TMUX" ] && [ -z "$GHOSTTY" ]; then
    if tmux has-session 2>/dev/null; then
        exec tmux attach
    else
        exec tmux new-session
    fi
fi

# 自动 tmux 分屏逻辑 (仅针对第一个窗口)
if [[ -n "$TMUX" && -z "$SPLITED" ]]; then
    # 获取当前窗口的索引号
    CURRENT_WINDOW=$(tmux display-message -p '#I')

    # 判断是否为第一个窗口 (通常是 0，如果你设置了 base-index 1 则改为 1)
    if [[ "$CURRENT_WINDOW" == "0" ]]; then
        # 检查当前窗口是否只有一个面板，防止重复触发
        if [ "$(tmux list-panes | wc -l)" -eq 1 ]; then
            # 标记已分屏，防止当前窗口的子 shell 再次触发
            export SPLITED=1
            
            # 执行分屏逻辑
            tmux split-window -h -p 30
            tmux split-window -v -p 50
            tmux select-pane -t 0
        fi
    fi
fi

function zvm_after_init() {
  bindkey "^R" fzf-history-widget
}

source $(brew --prefix)/opt/zsh-vi-mode/share/zsh-vi-mode/zsh-vi-mode.plugin.zsh
ZVM_VI_INSERT_ESCAPE_BINDKEY=jk

source $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh

bindkey -M emacs -r "^R"
bindkey -M viins -r "^R"
source <(fzf --zsh)

fpath+=("$(brew --prefix)/share/zsh/site-functions")
autoload -U promptinit; promptinit
# optionally define some options
PURE_CMD_MAX_EXEC_TIME=10
# change the path color
zstyle :prompt:pure:path color white
# change the color for both `prompt:success` and `prompt:error`
zstyle ':prompt:pure:prompt:*' color cyan
# turn on git stash status
zstyle :prompt:pure:git:stash show yes
prompt pure

