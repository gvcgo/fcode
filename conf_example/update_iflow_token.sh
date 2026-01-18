#!/bin/zsh

iflow xxxxxx > /dev/null 2>&1

iflowApiKey=$(python3 -c "import json; print(json.load(open('$HOME/.iflow/iflow_accounts.json'))['iflowApiKey'])")

mkdir -p "$HOME/.local/share/opencode"

auth_file="$HOME/.local/share/opencode/auth.json"

if [[ ! -f "$auth_file" ]]; then
    python3 -c "import json; json.dump({}, open('$auth_file', 'w'))"
fi

python3 -c "
import json

with open('$auth_file', 'r') as f:
    data = json.load(f)

if 'iflowcn' not in data:
    data['iflowcn'] = {'type': 'api'}

data['iflowcn']['key'] = '$iflowApiKey'

with open('$auth_file', 'w') as f:
    json.dump(data, f, indent=2)
"
