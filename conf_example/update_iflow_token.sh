#!/bin/zsh

# 脚本功能：
# 自动刷新 iflow 的 API token 并同步到 OpenCode 的认证文件
#
# 使用方法：
#   ./update_iflow_token.sh [选项]
#
# 选项：
#   -h, --help              显示帮助信息
#   -w, --wait-time SECONDS 等待时间（默认：5秒）
#   -i, --iflow-config FILE iflow 配置文件路径（默认：~/.iflow/iflow_accounts.json）
#   -a, --auth-file FILE    OpenCode 认证文件路径（默认：~/.local/share/opencode/auth.json）
#   -v, --verbose           显示详细输出
#   --dry-run               仅模拟执行，不实际修改文件

set -euo pipefail

# ==================== 配置 ====================

# 默认值
DEFAULT_WAIT_TIME=5
DEFAULT_IFLOW_ACCOUNTS_FILE="$HOME/.iflow/iflow_accounts.json"
DEFAULT_AUTH_FILE="$HOME/.local/share/opencode/auth.json"

# 变量（可通过命令行覆盖）
WAIT_TIME=$DEFAULT_WAIT_TIME
IFLOW_ACCOUNTS_FILE=$DEFAULT_IFLOW_ACCOUNTS_FILE
AUTH_FILE=$DEFAULT_AUTH_FILE
VERBOSE=false
DRY_RUN=false

# ==================== 函数 ====================

# 显示帮助信息
show_help() {
    cat << 'EOF'
用法: update_iflow_token.sh [选项]

自动刷新 iflow 的 API token 并同步到 OpenCode 的认证文件。

选项：
  -h, --help              显示此帮助信息
  -w, --wait-time SECONDS 等待 iflow 刷新 token 的时间（默认：5秒）
  -i, --iflow-config FILE iflow 配置文件路径
                          （默认：~/.iflow/iflow_accounts.json）
  -a, --auth-file FILE    OpenCode 认证文件路径
                          （默认：~/.local/share/opencode/auth.json）
  -v, --verbose           显示详细输出
  --dry-run               模拟执行，不实际修改文件

示例：
  ./update_iflow_token.sh                              # 使用默认配置
  ./update_iflow_token.sh -w 10 -v                    # 等待10秒，显示详细输出
  ./update_iflow_token.sh --iflow-config /path/to/json --auth-file /path/to/auth
  ./update_iflow_token.sh --dry-run                   # 查看将要执行的操作
EOF
}

# 日志函数
log_info() {
    echo "ℹ️  $*"
}

log_verbose() {
    if [[ "$VERBOSE" == true ]]; then
        echo "🔍 $*"
    fi
}

log_success() {
    echo "✅ $*"
}

log_warning() {
    echo "⚠️  $*" >&2
}

log_error() {
    echo "❌ $*" >&2
}

# 获取可用的 Python 命令
get_python_cmd() {
    if command -v python3 &>/dev/null; then
        echo "python3"
    elif command -v python &>/dev/null; then
        # 检查是否是 Python 2
        local python_version
        python_version=$(python --version 2>&1)
        if [[ "$python_version" =~ Python\ 2\.[0-9] ]]; then
            echo "python"
        else
            echo "python"
        fi
    else
        echo ""
    fi
}

# 从 JSON 文件读取值
read_json_value() {
    local file=$1
    local key=$2

    local python_cmd
    python_cmd=$(get_python_cmd)

    if [[ -z "$python_cmd" ]]; then
        log_error "需要 Python 来解析 JSON 文件"
        return 1
    fi

    log_verbose "使用 $python_cmd 读取 $file 的 $key"

    $python_cmd -c "import json; print(json.load(open('$file'))['$key'])" 2>/dev/null || {
        log_error "无法从 $file 读取 $key"
        return 1
    }
}

# 写入 JSON 文件
write_json_value() {
    local file=$1
    local section=$2
    local key=$3
    local value=$4

    if [[ "$DRY_RUN" == true ]]; then
        log_info "[DRY RUN] 将写入 $file: $section.$key = ${value:0:10}..."
        return 0
    fi

    local python_cmd
    python_cmd=$(get_python_cmd)

    if [[ -z "$python_cmd" ]]; then
        log_error "需要 Python 来更新 JSON 文件"
        return 1
    fi

    log_verbose "使用 $python_cmd 更新 $file"

    $python_cmd << EOF
import json

with open('$file', 'r') as f:
    data = json.load(f)

if '$section' not in data:
    data['$section'] = {'type': 'api'}

data['$section']['$key'] = '$value'

with open('$file', 'w') as f:
    json.dump(data, f, indent=2)

print("已更新: $file")
EOF
}

# 安全终止进程
kill_process_safely() {
    local pid=$1
    local timeout=${2:-5}

    log_verbose "尝试终止进程 $pid (超时: ${timeout}秒)..."

    # 先发送 SIGTERM
    kill "$pid" 2>/dev/null || return 0

    # 等待进程退出
    local count=0
    while kill -0 "$pid" 2>/dev/null; do
        sleep 0.5
        ((count++)) || true
        if ((count >= timeout * 2)); then
            log_warning "进程 $pid 在 ${timeout}秒内未响应，强制终止"
            kill -9 "$pid" 2>/dev/null
            break
        fi
    done

    # 等待僵尸进程
    wait "$pid" 2>/dev/null || true

    log_verbose "进程 $pid 已终止"
}

# ==================== 参数解析 ====================

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -w|--wait-time)
                if [[ -z "${2:-}" ]] || [[ "$2" =~ ^- ]]; then
                    log_error "选项 $1 需要一个参数"
                    exit 1
                fi
                if ! [[ "$2" =~ ^[0-9]+$ ]]; then
                    log_error "等待时间必须是正整数"
                    exit 1
                fi
                WAIT_TIME=$2
                shift 2
                ;;
            -i|--iflow-config)
                if [[ -z "${2:-}" ]] || [[ "$2" =~ ^- ]]; then
                    log_error "选项 $1 需要一个参数"
                    exit 1
                fi
                IFLOW_ACCOUNTS_FILE=$2
                shift 2
                ;;
            -a|--auth-file)
                if [[ -z "${2:-}" ]] || [[ "$2" =~ ^- ]]; then
                    log_error "选项 $1 需要一个参数"
                    exit 1
                fi
                AUTH_FILE=$2
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# ==================== 主流程 ====================

main() {
    parse_args "$@"

    if [[ "$DRY_RUN" == true ]]; then
        log_info "=== DRY RUN 模式 ==="
    fi

    log_verbose "配置:"
    log_verbose "  iflow 配置文件: $IFLOW_ACCOUNTS_FILE"
    log_verbose "  认证文件: $AUTH_FILE"
    log_verbose "  等待时间: ${WAIT_TIME}秒"

    # 步骤 1: 在后台启动 iflow
    log_info "步骤 1: 在后台启动 iflow..."

    if ! command -v iflow &>/dev/null; then
        log_error "未找到 iflow 命令，请确保已安装"
        exit 1
    fi

    if [[ "$DRY_RUN" == false ]]; then
        iflow xxxxxx &>/dev/null
        IFLOW_PID=$!
        log_info "iflow 进程 PID: $IFLOW_PID"
        log_verbose "等待 ${WAIT_TIME} 秒让 iflow 刷新 token..."

        sleep "$WAIT_TIME"

        # 步骤 2: 关闭 iflow 进程
        log_info "步骤 2: 关闭 iflow 进程..."
        kill_process_safely "$IFLOW_PID"
    else
        log_info "[DRY RUN] 跳过启动 iflow 进程"
    fi

    # 步骤 3: 读取 iflowApiKey
    log_info "步骤 3: 读取 iflowApiKey..."

    if [[ ! -f "$IFLOW_ACCOUNTS_FILE" ]]; then
        log_error "文件不存在: $IFLOW_ACCOUNTS_FILE"
        exit 1
    fi

    IFLOW_API_KEY=$(read_json_value "$IFLOW_ACCOUNTS_FILE" "iflowApiKey")

    if [[ -z "$IFLOW_API_KEY" ]]; then
        log_error "无法从 $IFLOW_ACCOUNTS_FILE 读取 iflowApiKey"
        exit 1
    fi

    log_success "成功读取 iflowApiKey: ${IFLOW_API_KEY:0:10}..."

    # 步骤 4: 更新 auth.json
    log_info "步骤 4: 更新 auth.json..."

    if [[ ! -f "$AUTH_FILE" ]]; then
        log_info "文件 $AUTH_FILE 不存在，创建新文件..."
        if [[ "$DRY_RUN" == false ]]; then
            mkdir -p "$(dirname "$AUTH_FILE")"
            echo "{}" > "$AUTH_FILE"
        else
            log_info "[DRY RUN] 将创建 $AUTH_FILE"
        fi
    fi

    if write_json_value "$AUTH_FILE" "iflowcn" "key" "$IFLOW_API_KEY"; then
        log_success "完成！"
    else
        log_error "更新 auth.json 失败"
        exit 1
    fi
}

# 执行主函数
main "$@"
