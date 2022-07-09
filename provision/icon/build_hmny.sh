#!/usr/bin/env -S bash -eET -o pipefail
args=("${@:1}")

docker_user="ubuntu"
docker_host="localhost"
docker_port="5000"
docker_registry="$docker_host"
[[ -z $docker_port ]] || docker_registry+=":$docker_port"

hmny_dir=$PWD
hmny_tmp_dir="$hmny_dir/_hmny"
hmny_src_dir="$hmny_dir/res"
root_dir="$hmny_dir/../.."

function repeat() {
    for i in $(seq 1 $2); do echo -n "$1"; done
}

function build_images() {

    image="${1:-}"
    repos_dir=$hmny_tmp_dir/repos
    mkdir -p $repos_dir

    function build_hmny() {
        cd $repos_dir
        docker <$hmny_src_dir/hmny.Dockerfile \
            build \
            --build-arg SHARDING_HOST="$docker_host" \
            -t $docker_registry/hmny:latest -
        cd $hmny_dir
    }

    case "$image" in
    hmny) build_hmny ;;
    *)
        build_hmny
        ;;
    esac
}
build_images "${args[@]}"
