#!/bin/bash

if [ $# -ne 1 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node Max>"
    echo ""

    exit 0
fi

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
NODE_MAX=$1


$PROJ_PATH/mykill

sleep 2

$PROJ_PATH/start-geth-all.sh $NODE_MAX

