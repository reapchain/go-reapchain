#!/bin/bash

if [ $# -ne 2 ]; then
	echo ""
	echo "Usage : $(basename "$0") <Node Max> <Json File>"
	echo "        <Json File> sample josn 파일명"
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
SAMPLE_JSON_FILE="$SETUP_INFO_PATH/$2"
OUT_JSON_FILE=$( echo "$SAMPLE_JSON_FILE" | sed "s/\.sample//" )

#echo "SAMPLE_JSON_FILE : $SAMPLE_JSON_FILE"
#echo "OUT_JSON_FILE : $OUT_JSON_FILE"

JQ_BIN="$PROJ_PATH/bin/jq"


#####################################################
# JSON Parsing
#----------------------------------------------------
func_JsonParsing()
{
	ret_val=`cat $SAMPLE_JSON_FILE | $JQ_BIN $1`
	ret_val=${ret_val//\"/} # 앞뒤의 " 제거
	echo $ret_val
}

RPC_ADDR=`func_JsonParsing ".arguments.ip"`
RPC_PORT=`func_JsonParsing ".arguments.rpc_port"`

#echo "RPC_ADDR = $RPC_ADDR"
#echo "RPC_PORT = $RPC_PORT"

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

rpc_port=`expr $RPC_PORT + $node_no`
http_connect="http://$RPC_ADDR:$rpc_port"

json_fmt="{\"jsonrpc\":\"2.0\",\"id\":0,\"method\":\"eth_accounts\",\"params\":[\"\"]}"
ret_json=`curl -X POST -m 5 --data "$json_fmt" $http_connect`

echo "--------------------------------------------------------------"
#echo "[$node_no] $json_fmt"
#echo "[$node_no] $ret_json"

g_ret_getAccount=`echo $ret_json | $JQ_BIN '.result'`

if [ "$g_ret_getAccount" == "null" ];
then
	return
fi

#echo "g_ret_getAccount = $g_ret_getAccount"
g_ret_getAccount=${g_ret_getAccount//\"/} # '"' 제거
g_ret_getAccount=${g_ret_getAccount//\[/} # '[' 제거
g_ret_getAccount=${g_ret_getAccount//\]/} # ']' 제거
g_ret_getAccount=${g_ret_getAccount//\,/} # ',' 제거
#echo "g_ret_getAccount = $g_ret_getAccount"
}

#####################################################
# 잔고 조회 : curl
#----------------------------------------------------
func_getBalance_curl()
{
addr=$1
node_no=$2

rpc_port=`expr $RPC_PORT + $node_no`
http_connect="http://$RPC_ADDR:$rpc_port"

json_fmt="{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"$addr\",\"latest\"],\"id\":0}"
ret_json=`curl -X POST -m 5 --data "$json_fmt" $http_connect`

#echo "[$node_no] $json_fmt"
#echo "[$node_no] $ret_json"

g_ret_getBalance=`echo $ret_json | $JQ_BIN '.result'`

if [ "$g_ret_getBalance" == "null" ];
then
	return
fi

g_ret_getBalance=${g_ret_getBalance//\"/} # 앞뒤의 " 제거
#echo "[$node_no] $addr : $(( 10#$g_ret_getBalance ))"
echo "[$node_no] $addr : $(printf "%40s" $g_ret_getBalance)   [$(date +"%Y-%m-%d %H:%M:%S.%3N")]"
}

#####################################################
# 잔고 조회 : geth
#----------------------------------------------------
func_getBalance_geth()
{
addr=$1
node_no=$2

rpc_port=`expr $RPC_PORT + $node_no`
http_connect="http://$RPC_ADDR:$rpc_port"

cmd="eth.getBalance(\"$addr\")"

result=`$BIN_PATH/geth attach rpc:$http_connect << EOF
$cmd
EOF`

g_ret_getBalance=`func_getValueParsing "$result"`
echo "[$node_no] $addr : $(printf "%40s" $g_ret_getBalance)   [$(date +"%Y-%m-%d %H:%M:%S.%3N")]"
}

#####################################################
# Account 조회하여 JSON_FILE 에 치환
#----------------------------------------------------
main_accounts()
{
loop_max=$1
json_file=$2

loop_cnt=1
while [ $loop_cnt -lt $loop_max ];
do
	func_getAccount_curl $loop_cnt

	if [ "$g_ret_getAccount" != "null" ];
	then
		update_from="from_$loop_cnt"
		update_to="to_$loop_cnt"
		for_cnt=0
		for i in $g_ret_getAccount;
		do
			for_cnt=$(( $for_cnt + 1 ))
			case $for_cnt in
				1) update_str="from_${loop_cnt}_addr" ;;
				2) update_str="to_${loop_cnt}_addr" ;;
			esac

			echo "update_str = $update_str"

			$( sed -i "s/$update_str/$i/g" $json_file )

			[[ $for_cnt -eq 2 ]] && break

		done
		#echo "$update_str = $g_ret_getAccount"
	fi

	loop_cnt=$(( $loop_cnt + 1 ))
done
}

#####################################################
# main
#####################################################

cp -p $SAMPLE_JSON_FILE $OUT_JSON_FILE

main_accounts $NODE_MAX $OUT_JSON_FILE > $TLOG_FILE_NAME

