#!/bin/bash

dir=$(pwd)
echo "$dir"
# shellcheck disable=SC2164
cd "$dir"/plugin/nube/system/ && rm system.so
# shellcheck disable=SC2164
cd "$dir"/plugin/nube/protocals/lorawan && rm lorawan.so
# shellcheck disable=SC2164
cd "$dir"/plugin/nube/protocals/modbus && rm modbus.so
# shellcheck disable=SC2164
cd "$dir"/plugin/nube/protocals/lora && rm lora.so
