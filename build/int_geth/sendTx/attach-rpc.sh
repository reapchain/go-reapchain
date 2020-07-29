#!/bin/bash

if [ $# -lt 1 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node No> <Target IP>"
    echo ""

    return
fi

NODE_NO=$1
RPC_IP=$2

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
BIN_PATH=$PROJ_PATH/bin
RPC_PORT=$(( 8540 + $NODE_NO ))

if [ "$RPC_IP" == "" ];
then
	RPC_IP=$(hostname -I | awk '{print $1}')
fi

HTTP_CONNECT="http://$RPC_IP:$RPC_PORT"
echo $HTTP_CONNECT

$BIN_PATH/geth attach rpc:$HTTP_CONNECT

