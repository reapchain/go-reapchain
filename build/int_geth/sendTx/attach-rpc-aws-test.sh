#!/bin/sh

if [ $# -ne 1 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node No>"
    echo ""

    return
fi

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
BIN_PATH=$PROJ_PATH/bin
RPC_PORT=$(( 8540 + $1 ))
#RPC_IP=localhost
RPC_IP=13.125.221.41

$BIN_PATH/geth attach rpc:http://$RPC_IP:$RPC_PORT

