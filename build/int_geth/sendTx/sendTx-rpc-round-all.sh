#!/bin/bash

if [ $# -lt 1 ]; then
	echo ""
	echo "Usage : $(basename "$0") <max_node_no> [round_count] [sleep] [Json_File]"
	echo "        <max_node_no> 전체 node 개수"
	echo "        [round_count] reound 반복 건수 (미 지정이나 0 지정 시 무한 반복)"
	echo "        [sleep]       reound Tx 송신 후 sleep time"
	echo "        [Json File]   Json config file (필수 입력 항목 아님)"
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
MAX_NODE_NO=$1
ROUND_COUNT=$2
SLEEP_TIME=$3
JSON_CONFIG_FILE=$4


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
TIMESTAMP_FILE=$TLOG_PATH/timestamp_${BASE_NAME}.log
TLOG_FILE_NAME="${TLOG_PATH}/${BASE_NAME}.log"

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
TX_TYPE=`func_JsonParsing ".arguments.tx_type"`


#####################################################
func_getValueParsing()
{
result_strings=$1
findok=false
for i in $result_strings ;
do
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
func_getCmdSaveFile()
{
	arg_from_node_no=$1
	arg_max_node_no=$2

	echo "$TMP_PATH/cmd_${BASE_NAME}_${arg_from_node_no}_to_${arg_max_node_no}.tmp"
}


#####################################################
# SendTransaction 생성
#----------------------------------------------------
func_createTx()
{
	arg_type=$1
	arg_from_node_no=$2
	arg_max_node=$3
	max_node_no=$((${arg_max_node}-1))

	from_node_no=`expr $arg_from_node_no + 0`

	from_addr=`func_JsonParsing ".node_${from_node_no}.${arg_type}.from"`
	send_value=`func_JsonParsing ".node_${from_node_no}.${arg_type}.value"`
	send_data=`func_JsonParsing ".node_${from_node_no}.${arg_type}.data"`

	#echo "func_createTx: arg_type     = $arg_type"
	#echo "func_createTx: from_node_no = $from_node_no"
	#echo "func_createTx: arg_max_node = $arg_max_node"
	#echo "func_createTx: send_value   = $send_value"
	#echo "func_createTx: send_data    = $send_data"
	#echo "func_createTx: from_addr    = $from_addr"

	# unlock Account
	cmd_unlock="personal.unlockAccount(\"$from_addr\", \"reapchain\", 0)"

	target_rpc_port=$(( ${RPC_PORT} + ${from_node_no} ))
	#echo "func_createTx: target_rpc_port = $target_rpc_port"

	func_attach "http://$RPC_IP:$target_rpc_port" "$cmd_unlock"

	cmd_save_file=$(func_getCmdSaveFile $arg_from_node_no $arg_max_node)
	$(>$cmd_save_file)

	#echo "func_createTx: cmd_save_file = $cmd_save_file"

	# sendTx command 생성
	to_node_no=1
	for to_node_no in `seq 1 $max_node_no`;
	do
		to_addr=`func_JsonParsing ".node_${to_node_no}.${arg_type}.to"`

		if [ "$arg_type" == "messages" ];
		then
			cmd_fmt="eth.sendTransaction({from: \"$from_addr\", to: \"$to_addr\", value: \"$send_value\"})"
		else
			cmd_fmt="eth.sendTransaction({from: \"$from_addr\", to: \"$to_addr\", value: \"$send_value\", data: \"$send_data\"})"
		fi

		echo "$cmd_fmt" >> $cmd_save_file
	done
}

#####################################################
# SendTransaction processing
#----------------------------------------------------
func_sendTx()
{
	arg_round_cnt=$1
	arg_sleep=$2
	arg_type=$3
	arg_max_node=$4
	last_node_no=$(($arg_max_node-1))

	#echo "func_sendTx: arg_round_cnt = $arg_round_cnt"
	#echo "func_sendTx: arg_sleep     = $arg_sleep"
	#echo "func_sendTx: arg_type      = $arg_type"
	#echo "func_sendTx: arg_max_node  = $arg_max_node"
	#echo "func_sendTx: last_node_no  = $last_node_no"


	# node_no 기준으로 unlock & sendTx command 생성
	for node_no in `seq -w 1 $last_node_no`;
	do
		func_createTx $arg_type $node_no $arg_max_node
	done


	# Tx 송신 처리

	loop_cnt=1
	if [ $arg_round_cnt -eq 0 ]; then
		loop_max=1
	else
		loop_max=$arg_round_cnt
	fi

	while [ $loop_cnt -le $loop_max ];
	do
		echo "[$loop_cnt] $(date +"%Y-%m-%d %H:%M:%S.%3N")" >> $TIMESTAMP_FILE

		for from_node_no in `seq -w 1 $last_node_no`;
		do
			cmd_save_file=$(func_getCmdSaveFile $from_node_no $arg_max_node)
			cmd=$(<$cmd_save_file)

			target_rpc_port=$(( ${RPC_PORT} + `expr $from_node_no + 0` ))
			#echo "func_sendTx: target_rpc_port = $target_rpc_port"
			echo "$cmd_save_file"

			func_attach "http://$RPC_IP:$target_rpc_port" "$cmd"

			if [ $arg_sleep -gt 0 ]; then
				sleep $arg_sleep
			fi
		done

		echo "[$loop_cnt] $(date +"%Y-%m-%d %H:%M:%S.%3N")" >> $TIMESTAMP_FILE

		loop_cnt=$(($loop_cnt + 1))
		if [ $arg_round_cnt -eq 0 ]; then
			loop_max=$(($loop_cnt + 1))
		fi
	done
}

#####################################################
# main
#####################################################

# 파일 초기화
rm -f $TMP_PATH/cmd_*
$(> $TIMESTAMP_FILE)

func_sendTx $ROUND_COUNT $SLEEP_TIME $TX_TYPE $MAX_NODE_NO > $TLOG_FILE_NAME

