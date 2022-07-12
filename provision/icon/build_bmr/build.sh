#!/usr/bin/env -S bash -eET -o pipefail
args=("${@:1}")

docker_user="ubuntu"
docker_host="localhost"
docker_port="5000"
docker_registry="$docker_host"
[[ -z $docker_port ]] || docker_registry+=":$docker_port"

bmr_dir=$PWD
bmr_tmp_dir="$bmr_dir/_bmr"
bmr_src_dir="$bmr_dir/res"
root_dir="$bmr_dir/../.."

function build_images() {
    image="${1:-}"
    repos_dir=$bmr_tmp_dir/repos
    mkdir -p $repos_dir

    function build_bmr() {
        cd $root_dir
        docker \
            build \
            -f $bmr_src_dir/Dockerfile \
            -t $docker_registry/bmr:latest .
        cd $bmr_dir
    }

    case "$image" in
    bmr) build_bmr ;;
    *)
        build_bmr
        ;;
    esac
}
build_images "${args[@]}"
