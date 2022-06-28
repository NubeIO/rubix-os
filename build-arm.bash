#!/bin/bash

# Console colors
DEFAULT="\033[0m"
RED="\033[31m"
YELLOW="\033[33m"
GREEN="\033[32m"

help() {
  echo "USAGE: see build.bash --help for same options:"
  echo ""
  . build.bash --help
}

parseCommand() {
  for i in "$@"; do
    case ${i} in
    -h | --help)
      help
      exit 0
      ;;
    esac
  done
}

parseCommand "$@"
set -e
docker build --file Dockerfile.armv7 --tag go-build-flow-framework-armv7 --build-arg plugins="$*" .
docker container create --name temp go-build-flow-framework-armv7
set +e
docker container cp temp:/app/flow-framework.armv7.zip ./
docker container rm temp

echo -e "${GREEN}OUTPUT: flow-framework.armv7.zip${DEFAULT}"