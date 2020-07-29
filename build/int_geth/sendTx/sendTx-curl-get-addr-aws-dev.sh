#!/bin/bash

PROJ_PATH="$( cd "$DIR_NAME" ; pwd -P )"

$PROJ_PATH/sendTx-curl-get-addr.sh 50 15.164.245.122 8640 sendTx-rpc.json.aws.dev

