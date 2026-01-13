function zvm_after_init() {
  bindkey "^R" fzf-history-widget
}

source $(brew --prefix)/opt/zsh-vi-mode/share/zsh-vi-mode/zsh-vi-mode.plugin.zsh
ZVM_VI_INSERT_ESCAPE_BINDKEY=jk

bindkey -M emacs -r "^R"
bindkey -M viins -r "^R"
source <(fzf --zsh)
