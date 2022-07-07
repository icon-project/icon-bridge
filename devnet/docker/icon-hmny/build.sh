#!/usr/bin/env -S bash -eET -o pipefail
args=("${@:2}")

docker_user="ubuntu"
docker_host="localhost"
docker_port="5000"
docker_registry="$docker_host"
[[ -z $docker_port ]] || docker_registry+=":$docker_port"

ixh_dir=$PWD
ixh_tmp_dir=$ixh_dir/_ixh
ixh_build_dir=$ixh_tmp_dir/build
ixh_src_dir=$ixh_dir/src

root_dir="$ixh_dir/../../.."

btp_icon_branch="v1.2.3"
btp_hmny_branch="v4.3.7"

# localnet: begin
btp_icon_config=$ixh_src_dir/icon.config.json

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
    function build_icon() {
        log "building icon"
        cd $repos_dir
        if [ -d goloop ]; then
            cd goloop
            git fetch
            git checkout ${btp_icon_branch:-master}
        else
            git clone --single-branch \
                --branch ${btp_icon_branch:-master} \
                https://github.com/icon-project/goloop
            cd goloop
        fi
    echo "here"
        make gochain-icon-image
        docker <$ixh_src_dir/icon.Dockerfile \
            build \
            --build-arg CONFIG_JSON="$(cat $btp_icon_config)" \
            -t $docker_registry/icon:latest -
        cd $ixh_dir
    echo "there"
    }
    function build_hmny() {
        log "building hmny"
        cd $repos_dir
        # git clone --single-branch \
        #     --branch ${btp_icon_branch:-main} \
        #     https://github.com/harmony-one/harmony
        # cd harmony
        docker <$ixh_src_dir/hmny.Dockerfile \
            build \
            --build-arg SHARDING_HOST="$docker_host" \
            -t $docker_registry/hmny:latest -
        cd $ixh_dir
    }

    case "$image" in
    bmr) build_bmr ;;
    icon) build_icon ;;
    hmny) build_hmny ;;
    *)
        build_hmny
        build_icon
        build_bmr
        ;;
    esac
}
build_images "${args[@]}"
