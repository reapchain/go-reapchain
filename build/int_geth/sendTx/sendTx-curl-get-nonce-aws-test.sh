#!/bin/bash

PROJ_PATH="$( cd "$DIR_NAME" ; pwd -P )"
PROC_TYPE=$1

$PROJ_PATH/sendTx-curl-get-nonce.sh 50 13.125.221.41 8540 sendTx-rpc.json.aws $PROC_TYPE

