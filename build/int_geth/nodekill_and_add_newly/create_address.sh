#!/bin/bash

if [ $# -ne 1 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node Count>"
    echo ""

    exit 0
fi

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
BIN_PATH=$PROJ_PATH/bin
NODE_NAME=node$1
DATA_DIR=$PROJ_PATH/$NODE_NAME

$BIN_PATH/geth --datadir $DATA_DIR account new --password $SETUP_INFO_PATH/passwd.txt
