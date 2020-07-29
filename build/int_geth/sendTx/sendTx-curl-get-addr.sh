#!/bin/bash

if [ $# -lt 1 ]; then
	echo ""
	echo "Usage : $(basename "$0") <Node Max> [Target IP] [RPC Port] [Json File]"
	echo "        <Node Max> 기동되어 있는 Node 개수"
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

NODE_MAX_LEN=${#NODE_MAX}

JQ_BIN="$PROJ_PATH/bin/jq"

if [ -z $JSON_FILE ];
then
	JSON_FILE="$SETUP_INFO_PATH/sendTx-rpc.json"
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
g_ret_getAccount=""
g_ret_isCreateBlock=""
g_ret_getBalance=""

g_connect_err=0

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
# Account 조회 : curl
#----------------------------------------------------
func_getAccount_curl()
{
node_no=$1
prt_node_no=$(printf "%${NODE_MAX_LEN}s" $node_no)

rpc_port=$(( $RPC_PORT + $node_no ))
http_connect="http://$RPC_IP:$rpc_port"

json_fmt="{\"jsonrpc\":\"2.0\",\"id\":0,\"method\":\"eth_accounts\",\"params\":[\"\"]}"
ret_json=`curl -X POST -m 5 --data "$json_fmt" $http_connect`

echo "--------------------------------------------------------------"
#echo "[$prt_node_no] $json_fmt"
#echo "[$prt_node_no] $ret_json"

g_ret_getAccount=`echo $ret_json | $JQ_BIN '.result'`

if [ "$g_ret_getAccount" == "null" ];
then
	return
fi

g_ret_getAccount=${g_ret_getAccount//\"/} # '"' 제거
g_ret_getAccount=${g_ret_getAccount//\[/} # '[' 제거
g_ret_getAccount=${g_ret_getAccount//\]/} # ']' 제거
g_ret_getAccount=${g_ret_getAccount//\,/} # ',' 제거
g_ret_getAccount=${g_ret_getAccount// /}  # ' ' 제거
#echo "g_ret_getAccount = $g_ret_getAccount"
}

#####################################################
# 잔고 조회 : curl
#----------------------------------------------------
func_getBalance_curl()
{
addr=$1
node_no=$2
prt_node_no=$(printf "%${NODE_MAX_LEN}s" $node_no)

rpc_port=$(( $RPC_PORT + $node_no ))
http_connect="http://$RPC_IP:$rpc_port"

json_fmt="{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"$addr\",\"latest\"],\"id\":0}"
ret_json=`curl -X POST -m 5 --data "$json_fmt" $http_connect`

#echo "[$prt_node_no] $json_fmt"
#echo "[$prt_node_no] $ret_json"

g_ret_getBalance=`echo $ret_json | $JQ_BIN '.result'`

if [ "$g_ret_getBalance" == "null" ];
then
	return
fi

g_ret_getBalance=${g_ret_getBalance//\"/}    # 앞뒤의 "  제거
g_ret_getBalance=${g_ret_getBalance#0x}      # 앞뒤의 0x 제거
g_ret_getBalance=$((16#${g_ret_getBalance})) # Hex -> Decimal
echo "[$prt_node_no] $addr : $(printf "%40s" $g_ret_getBalance)   [$(date +"%Y-%m-%d %H:%M:%S.%3N")]"
}

#####################################################
# 잔고 조회 : geth
#----------------------------------------------------
func_getBalance_geth()
{
addr=$1
node_no=$2
prt_node_no=$(printf "%${NODE_MAX_LEN}s" $node_no)

rpc_port=$(( $RPC_PORT + $node_no ))
http_connect="http://$RPC_IP:$rpc_port"

cmd="eth.getBalance(\"$addr\")"

result=`$BIN_PATH/geth attach rpc:$http_connect << EOF
$cmd
EOF`

g_ret_getBalance=`func_getValueParsing "$result"`
echo "[$prt_node_no] $addr : $(printf "%40s" $g_ret_getBalance)   [$(date +"%Y-%m-%d %H:%M:%S.%3N")]"
}

#####################################################
# Account별 잔고 조회
#----------------------------------------------------
main_accounts()
{
loop_max=$1

loop_cnt=1
while [ $loop_cnt -lt $loop_max ];
do
	func_getAccount_curl $loop_cnt

	if [ "$g_ret_getAccount" != "null" ];
	then
		for i in $g_ret_getAccount;
		do
			#func_getBalance_curl $i $loop_cnt
			func_getBalance_geth $i $loop_cnt
		done
	fi

	loop_cnt=$(( $loop_cnt + 1 ))
done
}

#####################################################
# main
#####################################################

main_accounts $NODE_MAX > $TLOG_FILE_NAME

