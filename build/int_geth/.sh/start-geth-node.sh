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
RPC_PORT=`expr 8640 \+ "$1" `

SERVER_IP="172.31.6.169"
#SERVER_IP="172.31.6.169"
BOOTNODE_PORT=30391
#BOOTNODE_ENODE="67cf6aaaf981cfcd8bf1a2d7d27b7c2ad41a1c722ce0793fad842d4d9197c04937da752a415061ca06895a5d1b9681734b0e75968498ad11366d88a01dd11e12"
BOOTNODE_ENODE="5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c"



# node 
if [ $1 -eq 0 ]; then
NODE_NAME=qman

#echo $PORT_NO $RPC_PORT
#echo $LOG_PATH/$NODE_NAME

nohup $BIN_PATH/geth \
	--networkid $NETWORK_ID \
	--port $PORT_NO \
	--datadir $PROJ_PATH/$NODE_NAME \
	--mine --minerthreads 1 \
	--syncmode full \
	--maxpeers 500 \
	--rpc \
	--rpcaddr $SERVER_IP \
	--rpcport $RPC_PORT \
	--rpccorsdomain "*" \
	--rpcapi="db,eth,net,web3,personal,miner,admin,debug,ssh,txpool,PoDC" \
	--verbosity 4 \
	--nat none \
	--bootnodes enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@172.31.6.169:30391 \
	2>&1 > $LOG_PATH/$NODE_NAME.log &

elif [ $1 -eq 5 ]; then

NODE_NAME=node$1

#echo $PORT_NO $RPC_PORT
#echo $LOG_PATH/$NODE_NAME

nohup $BIN_PATH/geth \
	--networkid $NETWORK_ID \
	--port $PORT_NO \
	--datadir $PROJ_PATH/$NODE_NAME \
	--mine --minerthreads 1 \
	--syncmode full \
	--maxpeers 500 \
	--unlock 0,1 \
	--password $SETUP_INFO_PATH/passwd.txt \
	--rpc \
	--rpcaddr $SERVER_IP \
	--rpcport $RPC_PORT \
	--rpccorsdomain "*" \
	--rpcapi="db,eth,net,web3,personal,miner,admin,debug,ssh,txpool,PoDC" \
	--verbosity 4 \
	--nat none \
	--bootnodes enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@172.31.6.169:30391 \
	2>&1 > $LOG_PATH/$NODE_NAME.log &

elif [ $1 -eq 6 ]; then

NODE_NAME=node$1

#echo $PORT_NO $RPC_PORT
#echo $LOG_PATH/$NODE_NAME

nohup $BIN_PATH/geth \
	--networkid $NETWORK_ID \
	--port $PORT_NO \
	--datadir $PROJ_PATH/$NODE_NAME \
	--mine --minerthreads 1 \
	--syncmode full \
	--maxpeers 500 \
	--unlock 0,1 \
	--password $SETUP_INFO_PATH/passwd.txt \
	--rpc \
	--rpcaddr $SERVER_IP \
	--rpcport $RPC_PORT \
	--rpccorsdomain "*" \
	--rpcapi="db,eth,net,web3,personal,miner,admin,debug,ssh,txpool,PoDC" \
	--verbosity 4 \
	--nat none \
	--bootnodes enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@172.31.6.169:30391 \
	--governance \
	2>&1 > $LOG_PATH/$NODE_NAME.log &
else

NODE_NAME=node$1

#echo $PORT_NO $RPC_PORT
#echo $LOG_PATH/$NODE_NAME

nohup $BIN_PATH/geth \
	--networkid $NETWORK_ID \
	--port $PORT_NO \
	--datadir $PROJ_PATH/$NODE_NAME \
	--mine --minerthreads 1 \
	--syncmode full \
	--maxpeers 500 \
	--unlock 0,1 \
	--password $SETUP_INFO_PATH/passwd.txt \
	--rpc \
	--rpcaddr $SERVER_IP \
	--rpcport $RPC_PORT \
	--rpccorsdomain "*" \
	--rpcapi="db,eth,net,web3,personal,miner,admin,debug,ssh,txpool,PoDC" \
	--verbosity 4 \
	--nat none \
	--bootnodes enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@172.31.6.169:30391 \
	2>&1 > $LOG_PATH/$NODE_NAME.log &
fi
#	--unlock 0 \
#	--password $SETUP_INFO_PATH/passwd.txt \
