#!/bin/bash

if [ $# -ne 1 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node No>"
    echo ""

    return
fi

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
BIN_PATH=$PROJ_PATH/bin
LOG_PATH=$PROJ_PATH/log

NETWORK_ID=2017
PORT_NO=`expr 30500 \+ "$1" `
RPC_PORT=`expr 8540 \+ "$1" `

BOOTNODE_IP="192.168.0.2"
BOOTNODE_PORT=30301
#BOOTNODE_ENODE="67cf6aaaf981cfcd8bf1a2d7d27b7c2ad41a1c722ce0793fad842d4d9197c04937da752a415061ca06895a5d1b9681734b0e75968498ad11366d88a01dd11e12"
BOOTNODE_ENODE="5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c"



# node 
if [ $1 -eq 0 ]; then
NODE_NAME=qman

#echo $PORT_NO $RPC_PORT
#echo $LOG_PATH/$NODE_NAME 
$BIN_PATH/geth \
	--networkid $NETWORK_ID \
	--port $PORT_NO \
	--datadir $PROJ_PATH/$NODE_NAME/data \
	--mine --minerthreads 1 \
	--syncmode "full" \
	--rpc \
	--rpcport $RPC_PORT \
	--rpccorsdomain "*" \
	--rpcapi="db,eth,net,web3,personal,miner,admin,debug,ssh,txpool,PoDC" \
	--unlock 0 \
	--password $SETUP_INFO_PATH/passwd.txt \
	--verbosity 4 \
	--nat none \
	--bootnodes enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@192.168.0.2:30301 \
	2>> $LOG_PATH/$NODE_NAME.log &
 
else

NODE_NAME=node$1

#echo $PORT_NO $RPC_PORT
#echo $LOG_PATH/$NODE_NAME

$BIN_PATH/geth \
	--networkid $NETWORK_ID \
	--port $PORT_NO \
	--datadir $PROJ_PATH/$NODE_NAME/data \
	--mine --minerthreads 1 \
	--syncmode "full" \
	--rpc \
	--rpcaddr $BOOTNODE_IP \
	--rpcport $RPC_PORT \
	--rpccorsdomain "*" \
	--rpcapi="db,eth,net,web3,personal,miner,admin,debug,ssh,txpool,PoDC" \
	--unlock 0 \
	--password $SETUP_INFO_PATH/passwd.txt \
	--verbosity 4 \
	--nat none \
	--bootnodes enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@192.168.0.2:30301 \
	2>> $LOG_PATH/$NODE_NAME.log &
fi
