#!/bin/bash

PROJ_PATH="$( cd "$DIR_NAME" ; pwd -P )"
PROC_TYPE=$1

$PROJ_PATH/sendTx-curl-get-nonce.sh 50 15.164.245.122 8640 sendTx-rpc.json.aws.dev $PROC_TYPE

