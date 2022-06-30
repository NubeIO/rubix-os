#!/bin/bash

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
  pluginDir=/data/flow-framework/data/plugins
else
  echo -e "${YELLOW}We are running in development mode!${DEFAULT}"
fi

echo -e "Creating a plugin directory if does not exist at: ${pluginDir}"
rm -rf $pluginDir/* || true
mkdir -p $pluginDir

BUILD_ERROR=false

function buildPlugin {
  echo -e "${DEFAULT}BUILDING $1..."
  go build -buildmode=plugin -o $1.so $2/*.go && cp $1.so $pluginDir
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}BUILD ${DEFAULT}$1"
  else
    echo -e "${RED}ERROR BUILD ${DEFAULT}$1"
    BUILD_ERROR=true
  fi
}

pushd $dir > /dev/null

for i in "$@"; do
  case ${i} in
    system)
      buildPlugin "system" plugin/nube/system ;;
    edge28)
      buildPlugin "edge28" plugin/nube/protocals/edge28 ;;
    modbus)
      buildPlugin "modbus" plugin/nube/protocals/modbus ;;
    lora)
      buildPlugin "lora" plugin/nube/protocals/lora ;;
    bacnet)
      buildPlugin "bacnet" plugin/nube/protocals/bacnetserver ;;
    lorawan)
      buildPlugin "lorawan" plugin/nube/protocals/lorawan ;;
    bacnet_master)
      buildPlugin "bacnet_master" plugin/nube/protocals/bacnetmaster ;;
    history)
      buildPlugin "history" plugin/nube/database/history ;;
    influx)
      buildPlugin "influx" plugin/nube/database/influx ;;
    rubixio)
      buildPlugin "rubixio" plugin/nube/protocals/rubixio ;;
    modbusserver)
      buildPlugin "modbusserver" plugin/nube/protocals/modbusserver ;;
    postgres)
      buildPlugin "postgres" plugin/nube/database/postgres ;;
  esac
done

if [ ${BUILD_ERROR} == true ]; then
    exit -1
fi

popd > /dev/null

if [ ${BUILD_ONLY} == true ]; then
    exit 0
fi

if [ ${PRODUCTION} == true ]; then
  go run app.go -g /data/flow-framework -d data --prod
else
  go run app.go --auth=false
fi
