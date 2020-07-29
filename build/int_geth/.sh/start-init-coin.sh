#!/bin/bash

if [ $# -ne 2 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node Max> <Genesis File>"
    echo ""

    exit 0
fi

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
NODEKEY_PATH=$SETUP_INFO_PATH/nodekey
BIN_PATH=$PROJ_PATH/bin
PASSWD=$SETUP_INFO_PATH/passwd.txt

NODE_MAX=$1
GENESIS_FILE=$SETUP_INFO_PATH/$2

rm -rf $PROJ_PATH/log/*
rm -rf $PROJ_PATH/nohup.out
rm -rf $PROJ_PATH/level

## node 0~6
for ((i=0;i<$NODE_MAX;i++));
do

# qmanager
if [ $i -eq 0 ]; then
NODE_NAME=qman
NODEKEY_FILE="$NODEKEY_PATH/nodekey.$NODE_NAME"
else
NODE_NAME=node$i
NODEKEY_FILE="$NODEKEY_PATH/nodekey$i"
fi

DATA_DIR=$PROJ_PATH/$NODE_NAME
mkdir -p $DATA_DIR/geth

# Copy node information
cp -p $SETUP_INFO_PATH/qmanager-nodes.json $DATA_DIR
#cp -p $SETUP_INFO_PATH/static-nodes.json   $DATA_DIR

# nodekey copy
cp -p $NODEKEY_FILE $DATA_DIR/geth/nodekey

# Create a key for the account
$BIN_PATH/geth --datadir $DATA_DIR account import --password $PASSWD $DATA_DIR/geth/nodekey

# Initialization genesis.json
$BIN_PATH/geth --datadir $DATA_DIR init $GENESIS_FILE

done

