#!/bin/bash

PROJ_PATH="$( cd "$DIR_NAME" ; pwd -P )"

$PROJ_PATH/sendTx-curl-get-addr.sh 50 192.168.0.80 8640 sendTx-rpc.json

