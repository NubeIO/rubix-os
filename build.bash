#!/bin/bash

if [ -z ${GO_PATH+x} ]; then
    GO_PATH=go
fi

# Console colors
DEFAULT="\033[0m"
RED="\033[31m"
YELLOW="\033[33m"
GREEN="\033[32m"

PRODUCTION=false
BUILD_ONLY=false

help() {
  echo "USAGE: bash build.bash [OPTIONS] [PLUGINS...]"
  echo "  i.e. bash build.bash --production modbus lorawan"
  echo ""
  echo "Options:"
  echo " -h  --help :          Print this help"
  echo " --prod  --production: Add these suffix to start production"
  echo " --build-only :        Don't run program"
}

parseCommand() {
  for i in "$@"; do
    case ${i} in
    -h | --help)
      help
      exit 0
      ;;
    --prod | --production)
      PRODUCTION=true
      ;;
    --build-only)
      BUILD_ONLY=true
      ;;
    esac
  done
}

parseCommand "$@"

dir=$(pwd)
echo -e "Current working directory is: $dir${DEFAULT}"
pluginDir=$dir/data/plugins

if [ ${PRODUCTION} == true ]; then
  echo -e "${YELLOW}We are running in production mode!${DEFAULT}"
  pluginDir=/data/rubix-os/data/plugins
else
  echo -e "${YELLOW}We are running in development mode!${DEFAULT}"
fi

echo -e "Creating a plugin directory if does not exist at: ${pluginDir}"
rm -rf $pluginDir/* || true
mkdir -p $pluginDir

BUILD_ERROR=false

function buildPlugin() {
  echo -e "${DEFAULT}BUILDING $1..."
  $GO_PATH build -buildmode=plugin -o $1.so $2/*.go && cp $1.so $pluginDir
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}BUILT ${DEFAULT}$1"
  else
    echo -e "${RED}ERROR BUILDING ${DEFAULT}$1"
    BUILD_ERROR=true
  fi
}

pushd $dir >/dev/null

for i in "$@"; do
  case ${i} in
  system)
    buildPlugin "system" plugin/nube/system
    ;;
  edge28)
    buildPlugin "edge28" plugin/nube/protocals/edge28
    ;;
  modbus)
    buildPlugin "modbus" plugin/nube/protocals/modbus
    ;;
  lora)
    buildPlugin "lora" plugin/nube/protocals/lora
    ;;
  bacnetserver)
    buildPlugin "bacnetserver" plugin/nube/protocals/bacnetserver
    ;;
  lorawan)
    buildPlugin "lorawan" plugin/nube/protocals/lorawan
    ;;
  bacnetmaster)
    buildPlugin "bacnetmaster" plugin/nube/protocals/bacnetmaster
    ;;
  history)
    buildPlugin "history" plugin/nube/database/history
    ;;
  influx)
    buildPlugin "influx" plugin/nube/database/influx
    ;;
  edgeinflux)
    buildPlugin "edgeinflux" plugin/nube/database/edgeinflux
    ;;
  edgeazure)
    buildPlugin "edgeazure" plugin/nube/database/edgeazure
    ;;
  rubixio)
    buildPlugin "rubixio" plugin/nube/protocals/rubixio
    ;;
  modbusserver)
    buildPlugin "modbusserver" plugin/nube/protocals/modbusserver
    ;;
  rubixpointsync)
    buildPlugin "rubixpointsync" plugin/nube/protocals/rubixpointsync
    ;;
  postgres)
    buildPlugin "postgres" plugin/nube/database/postgres
    ;;
  networklinker)
    buildPlugin "networklinker" plugin/nube/protocals/networklinker
    ;;
  galvintmv)
    buildPlugin "galvintmv" plugin/nube/projects/galvintmv
    ;;
  thresholdalerts)
    buildPlugin "thresholdalerts" plugin/nube/projects/thresholdalerts
    ;;
  flatlinealerts)
    buildPlugin "flatlinealerts" plugin/nube/projects/flatlinealerts
    ;;
  statusmismatchalerts)
    buildPlugin "statusmismatchalerts" plugin/nube/projects/statusmismatchalerts
    ;;
  maplora)
    buildPlugin "maplora" plugin/nube/protocals/maplora
    ;;
  mapmodbus)
    buildPlugin "mapmodbus" plugin/nube/protocals/mapmodbus
    ;;
  inauroazuresync)
    buildPlugin "inauroazuresync" plugin/nube/database/inauroazuresync
    ;;
  cpsprocessing)
    buildPlugin "cpsprocessing" plugin/nube/projects/cpsprocessing
    ;;
  esac
done

if [ ${BUILD_ERROR} == true ]; then
  exit 1
fi

popd >/dev/null

if [ ${BUILD_ONLY} == true ]; then
  echo -e "${DEFAULT}BUILDING app"
  go build app.go
  echo -e "${GREEN}BUILT ${DEFAULT}app"
else
    if [ ${PRODUCTION} == true ]; then
      go run app.go -g /data/rubix-os -d data --prod
    else
      go run app.go --auth=false
    fi
fi