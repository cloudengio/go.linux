#!/usr/bin/env bash

root="$(pwd)"

if [[ ! -f "${root}/.multimod.yml" && ! -f "${root}/go.mod" ]]; then
     echo "no .multimod.yml or go.mod in ${root}"
     exit 1
fi

docker run -it --rm \
    -v"${root}":/usr/src/repo -w /usr/src/repo \
     --security-opt seccomp=unconfined go.linux \
     "$@"

