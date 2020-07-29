#!/bin/bash

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"

$PROJ_PATH/copy_bin.sh geth
$PROJ_PATH/copy_bin.sh bootnode
$PROJ_PATH/copy_bin.sh istanbul
