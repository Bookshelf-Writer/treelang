#!/bin/bash
# Original source: https://github.com/Bookshelf-Writer/scripts-for-integration/blob/main/_run/push-hook.sh
echo "[HOOK]" "Push"

run_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
values_dir="$run_dir/values"
script_dir="$run_dir/scripts"
root_path=$(cd "$run_dir/.." && pwd)

#############################################################################

go mod tidy
#go mod vendor

bash "$run_dir/creator_dependencies_Go.sh"
bash "$run_dir/creator_const_Go.sh"


#############################################################################
exit 0

