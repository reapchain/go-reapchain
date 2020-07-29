#!/bin/bash

PROJ_PATH="$( cd "$DIR_NAME" ; pwd -P )"

$PROJ_PATH/sendTx-curl-get-addr.sh 50 13.125.221.41 8540 sendTx-rpc.json.aws

