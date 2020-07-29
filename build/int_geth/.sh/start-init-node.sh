#!/bin/bash

if [ $# -ne 1 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node Max>"
    echo ""

    exit 0
fi

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
NODEKEY_PATH=$SETUP_INFO_PATH/nodekey
BIN_PATH=$PROJ_PATH/bin
PASSWD=$SETUP_INFO_PATH/passwd.txt
NODE_MAX=$1

rm -rf $PROJ_PATH/node*
rm -rf $PROJ_PATH/qman*
rm -rf $PROJ_PATH/log/*
rm -rf $PROJ_PATH/nohup.out
rm -rf $PROJ_PATH/level

## node 0~6
for ((i=0;i<$NODE_MAX;i++));
do

# qmanager
if [ $i -eq 0 ]; then
	NODE_NAME=qman
	NODEKEY_FILE="nodekey.$NODE_NAME"
else
	NODE_NAME=node$i
	NODEKEY_FILE="nodekey$i"
fi

DATA_DIR=$PROJ_PATH/$NODE_NAME
mkdir -p $DATA_DIR/geth

cp $SETUP_INFO_PATH/genesis.json          $DATA_DIR
cp $SETUP_INFO_PATH/qmanager-nodes.json   $DATA_DIR
cp $NODEKEY_PATH/$NODEKEY_FILE            $DATA_DIR/geth/nodekey

cp -p $SETUP_INFO_PATH/static-nodes.json   $PROJ_PATH/$NODE_NAME

# import  a key for the account
$BIN_PATH/geth --datadir $DATA_DIR account import --password $PASSWD $DATA_DIR/geth/nodekey

# Initialization genesis.json
# $BIN_PATH/geth init $SETUP_INFO_PATH/genesis.json --datadir $PROJ_PATH/$NODE_NAME
 # Initialization genesis.json
$BIN_PATH/geth --datadir $DATA_DIR init  $DATA_DIR/genesis.json 

#$BIN_PATH/geth --datadir $PROJ_PATH/$NODE_NAME account import --password $SETUP_INFO_PATH/passwd.txt $PROJ_PATH/$NODE_NAME/geth/nodekey

done

