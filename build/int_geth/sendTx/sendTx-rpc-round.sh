#!/bin/bash

if [ $# -lt 2 ]; then
	echo ""
	echo "Usage : $(basename "$0") <from_node_no> <max_node_no> [round_count] [sleep] [Json_File]"
	echo "        <from_node_no> Tx 송신 node 번호 (from account[0])"
	echo "        <max_node_no>  Tx 수신 전체 node 개수 (to account[1])"
	echo "        [round_count]  reound 반복 건수 (미 지정이나 0 지정 시 무한 반복)"
	echo "        [sleep] reound Tx 송신 후 sleep time"
	echo "        [Json File]    Json config file (필수 입력 항목 아님)"
	echo ""

	exit
fi

BASE_NAME="$(basename "$0")"
BASE_NAME=${BASE_NAME/\.sh/}
DIR_NAME=$(dirname "$0")
PROJ_PATH="$( cd "$DIR_NAME" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
BIN_PATH=$PROJ_PATH/bin
LOG_PATH=$PROJ_PATH/log
TLOG_PATH=$PROJ_PATH/tlog
TMP_PATH=$PROJ_PATH/tmp

# Arguments
FROM_NODE_NO=$1
MAX_NODE_NO=$2
ROUND_COUNT=$3
SLEEP_TIME=$4
JSON_CONFIG_FILE=$5


JQ_BIN="$PROJ_PATH/bin/jq"

# 무한 반복
if [ -z $ROUND_COUNT ]; then
	ROUND_COUNT=0
fi

if [ -z $SLEEP_TIME ]; then
	SLEEP_TIME=0
fi

if [ -z $JSON_CONFIG_FILE ]; then
	JSON_CONFIG_FILE="sendTx-rpc.json"
fi

JSON_FILE="$SETUP_INFO_PATH/$JSON_CONFIG_FILE"
TIMESTAMP_FILE=$TLOG_PATH/timestamp_${BASE_NAME}_${FROM_NODE_NO}_to_${MAX_NODE_NO}.log
CMD_SAVE_FILE=$TMP_PATH/cmd_${BASE_NAME}_${FROM_NODE_NO}_to_${MAX_NODE_NO}.tmp
TLOG_FILE_NAME=${TLOG_PATH}/${BASE_NAME}_${FROM_NODE_NO}_to_${MAX_NODE_NO}.log

#####################################################
# Decimal to Hexadecimal converter
#----------------------------------------------------
func_Decimal2Hex()
{
	hex_value=`printf "0x%x" $(($1))`
	echo $hex_value
}
#####################################################
# JSON Parsing
#----------------------------------------------------
func_JsonParsing()
{
	ret_val=`cat $JSON_FILE | $JQ_BIN $1`
	ret_val=${ret_val//\"/} # 앞뒤의 " 제거
	echo $ret_val
}

#RPC_IP=$(ifconfig enp0s3 | grep inet | grep netmask | awk {'print $2'})
#RPC_IP=$(hostname -I | awk '{print $1}')
#RPC_IP=192.168.35.41
#RPC_IP=localhost
RPC_IP=`func_JsonParsing ".arguments.ip"`
RPC_PORT=`func_JsonParsing ".arguments.rpc_port"`
RPC_PORT=$(( ${RPC_PORT} + ${FROM_NODE_NO} ))
HTTP_CONNECT="http://$RPC_IP:$RPC_PORT"
TX_TYPE=`func_JsonParsing ".arguments.tx_type"`

TX_FROM=`func_JsonParsing ".node_${FROM_NODE_NO}.${TX_TYPE}.from"`
TX_GAS=`func_JsonParsing ".node_${FROM_NODE_NO}.${TX_TYPE}.gas"`
TX_GASPRICE=`func_JsonParsing ".node_${FROM_NODE_NO}.${TX_TYPE}.gasPrice"`
TX_VALUE=`func_JsonParsing ".node_${FROM_NODE_NO}.${TX_TYPE}.value"`
TX_DATA=`func_JsonParsing ".node_${FROM_NODE_NO}.${TX_TYPE}.data"`

# Decimal -> Hex
TX_GAS=`func_Decimal2Hex $TX_GAS`
TX_GASPRICE=`func_Decimal2Hex $TX_GASPRICE`
TX_VALUE=`func_Decimal2Hex $TX_VALUE`


#####################################################
func_getValueParsing()
{
	result_strings=$1
	findok=false
	for i in $result_strings ; do
		if [ "$i" == ">" ];
		then
			findok=true
		else
			if [ $findok == true ];
			then
				echo $i
				return
			fi
			findok=false
		fi
	done
}

#####################################################
func_attach()
{
$BIN_PATH/geth attach rpc:$1 << EOF
$2
EOF
}

#####################################################
# SendTransaction processing
#----------------------------------------------------
func_sendTx()
{
	arg_round_cnt=$1
	arg_sleep=$2
	arg_type=$3
	arg_from=$4
	arg_max_node=$5
	arg_value=$6
	arg_data=$7

	echo "==============================================================="

	# unlock Account
	cmd_unlock="personal.unlockAccount(\"$arg_from\", \"reapchain\", 0)"

	func_attach "$HTTP_CONNECT" "$cmd_unlock"

	# Tx 송신
	node_no=1
	while [ $node_no -lt $arg_max_node ];
	do
		to_addr=`func_JsonParsing ".node_${node_no}.${arg_type}.to"`

		if [ "$arg_type" == "messages" ];
		then
			cmd_fmt="eth.sendTransaction({from: \"$arg_from\", to: \"$to_addr\", value: \"$arg_value\"})"
		else
			cmd_fmt="eth.sendTransaction({from: \"$arg_from\", to: \"$to_addr\", value: \"$arg_value\", data: \"$arg_data\"})"
		fi

		echo "$cmd_fmt" >> $CMD_SAVE_FILE

		node_no=$(( $node_no + 1 ))
	done

	cmd=$(<$CMD_SAVE_FILE)


	loop_cnt=0
	if [ $arg_round_cnt -eq 0 ]; then
		loop_max=1
	else
		loop_max=$arg_round_cnt
	fi

	while [ $loop_cnt -lt $loop_max ];
	do

		loop_cnt=$(($loop_cnt + 1))
		if [ $arg_round_cnt -eq 0 ]; then
			loop_max=$(($loop_cnt + 1))
		fi

		echo "[$loop_cnt] $(date +"%Y-%m-%d %H:%M:%S.%3N")" >> $TIMESTAMP_FILE

		func_attach "$HTTP_CONNECT" "$cmd"

		echo "[$loop_cnt] $(date +"%Y-%m-%d %H:%M:%S.%3N")" >> $TIMESTAMP_FILE

		if [ $arg_sleep -gt 0 ]; then
			sleep $arg_sleep
		fi

	done

	rm -f $CMD_SAVE_FILE
}

#####################################################
# main
#####################################################

# 파일 초기화
$(> $CMD_SAVE_FILE)
$(> $TIMESTAMP_FILE)

func_sendTx $ROUND_COUNT $SLEEP_TIME $TX_TYPE $TX_FROM $MAX_NODE_NO $TX_VALUE $TX_DATA > $TLOG_FILE_NAME

