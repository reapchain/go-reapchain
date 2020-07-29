#!/bin/bash

if [ $# -ne 1 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node Max>"
    echo ""

    exit 0
fi

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
NODE_MAX=$1

## node create, make directory, import account, geth init 
$PROJ_PATH/start-bootnode.sh

## node startup
for ((i=0;i<$NODE_MAX;i++));
do
$PROJ_PATH/start-geth-node.sh $i
done

