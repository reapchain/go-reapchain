#!/bin/bash

if [ $# -lt 1 ]; then
	echo ""
	echo "Usage : $(basename "$0") <Node Max> [Target IP] [RPC Port] [Json File] [Process Type]"
	echo "        <Node Max> 기동되어 있는 Node 개수"
	echo "        [Process Type] latest, pending, earliest (default: latest)"
	echo "        [RPC Port] PRC 접속포트 (AWS-DEV:8640, AWS-TEST:8540)"
	echo ""

	exit
fi

BASE_NAME="$(basename "$0")"
DIR_NAME=$(dirname "$0")
PROJ_PATH="$( cd "$DIR_NAME" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
BIN_PATH=$PROJ_PATH/bin
LOG_PATH=$PROJ_PATH/log
TLOG_PATH=$PROJ_PATH/tlog

# Arguments
NODE_MAX=$1
RPC_IP=$2
RPC_PORT=$3
JSON_FILE=$4
PROC_TYPE=$5

NODE_MAX_LEN=${#NODE_MAX}

JQ_BIN="$PROJ_PATH/bin/jq"

if [ -z $JSON_FILE ];
then
	JSON_FILE="$SETUP_INFO_PATH/sendTx-rpc.json"
else
	JSON_FILE="$SETUP_INFO_PATH/$JSON_FILE"
fi

if [ -z $PROC_TYPE ];
then
	PROC_TYPE="latest"
fi
#####################################################
# JSON Parsing
#----------------------------------------------------
func_JsonParsing()
{
	ret_val=`cat $JSON_FILE | $JQ_BIN $1`
	ret_val=${ret_val//\"/} # 앞뒤의 " 제거
	echo $ret_val
}

if [ -z $RPC_IP ];
then
	RPC_IP=`func_JsonParsing ".arguments.ip"`
fi

if [ -z $RPC_PORT ];
then
	RPC_PORT=`func_JsonParsing ".arguments.rpc_port"`
fi

TLOG_FILE_NAME="${TLOG_PATH}/${BASE_NAME/\.sh/}.log"


#####################################################
# Shared variable
#----------------------------------------------------
g_ret_getNonce=""


#####################################################
# 계정 Nonce 조회
#----------------------------------------------------
func_getTransactionCount_curl()
{
	addr=$1
	node_no=$2
	proc_type=$3

	prt_node_no=$(printf "%${NODE_MAX_LEN}s" $node_no)

	rpc_port=$(( $RPC_PORT + $node_no ))
	http_connect="http://$RPC_IP:$rpc_port"

	#----------------------------------------------
	# "earliest" for the earliest/genesis block
	# "latest"   for the latest mined block
	# "pending"  for the pending state/transactions
	#----------------------------------------------
	json_fmt="{\"jsonrpc\":\"2.0\",\"id\":0,\"method\":\"eth_getTransactionCount\",\"params\":[\"$addr\",\"$proc_type\"]}"
	ret_json=`curl -X POST -H "Content-Type: application/json" -m 5 --data "$json_fmt" $http_connect`

	#echo "[$prt_node_no] $json_fmt"
	#echo "[$prt_node_no] $ret_json"

	g_ret_getNonce=`echo $ret_json | $JQ_BIN '.result'`

	if [ "$g_ret_getNonce" == "null" ];
	then
		return
	fi

	g_ret_getNonce=${g_ret_getNonce//\"/}    # 앞뒤의 "  제거
	g_ret_getNonce=${g_ret_getNonce#0x}      # 앞뒤의 0x 제거
	g_ret_getNonce=$((16#${g_ret_getNonce})) # Hex -> Decimal
	echo "-----------------------------------------------------"
	echo "[$prt_node_no] $addr : $(printf "%10s" $g_ret_getNonce)   [$(date +"%Y-%m-%d %H:%M:%S.%3N")]"
}

#####################################################
# Account별 Nonce 조회
#----------------------------------------------------
main_accounts()
{
	node_max=$1
	proc_type=$2

	tx_total_count=0
	node_no=1
	while [ $node_no -lt $node_max ];
	do

		from_addr=`func_JsonParsing ".node_${node_no}.messages.from"`
		#echo "from_addr = $from_addr"
		func_getTransactionCount_curl $from_addr $node_no $proc_type
		tx_total_count=$(( $tx_total_count + $g_ret_getNonce ))
		node_no=$(( $node_no + 1 ))
	done

	echo "-----------------------------------------------------"
	echo "($PROC_TYPE) [Total] Tx Count : $tx_total_count"
}

#####################################################
# main
#####################################################

main_accounts $NODE_MAX $PROC_TYPE > $TLOG_FILE_NAME

