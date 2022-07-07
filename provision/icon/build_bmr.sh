#!/usr/bin/env -S bash -eET -o pipefail
args=("${@:1}")

docker_user="ubuntu"
docker_host="localhost"
docker_port="5000"
docker_registry="$docker_host"
[[ -z $docker_port ]] || docker_registry+=":$docker_port"

ixh_dir=$PWD
ixh_tmp_dir=$ixh_dir/_ixh
ixh_src_dir=$ixh_dir/../../devnet/docker/icon-hmny/src
root_dir="$ixh_dir/../../"

function repeat() {
    for i in $(seq 1 $2); do echo -n "$1"; done
}

# echo message to stderr
function log() {
    local prefix="$(date '+%Y-%m-%d %H:%M:%S') $(repeat '    ' $((${#FUNCNAME[@]} - 2)))"
    echo -e "$prefix$@" >&2
}

function log_status() {
    [[ "$1" == 0 ]] && log " ✔" || log " ✘"
}

function log_stack() {
    local cmd=${FUNCNAME[1]}
    if [[ $# > 0 ]]; then cmd="$@"; fi
    local prefix="$(date '+%Y-%m-%d %H:%M:%S') $(repeat '    ' $((${#FUNCNAME[@]} - 3)))"
    echo -e "$prefix$cmd():${BASH_LINENO[1]}" >&2
}

function build_images() {
    log_stack

    image="${1:-}"
    repos_dir=$ixh_tmp_dir/repos
    mkdir -p $repos_dir

    function build_bmr() {
        log "building bmr"
        cd $root_dir
        docker \
            build \
            -f $ixh_src_dir/bmr.Dockerfile \
            -t $docker_registry/bmr:latest .
        cd $ixh_dir
    }

    case "$image" in
    bmr) build_bmr ;;
    *)
        build_bmr
        ;;
    esac
}
build_images "${args[@]}"
