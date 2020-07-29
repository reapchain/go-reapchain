#!/bin/bash

if [ $# -lt 3 ]; then
	echo ""
	echo "Usage : $(basename "$0") <Node No> <Send Count> <Sequence> [Json File]"
	echo "        <Send Count> Tx 전송 횟수"
	echo "        <Sequence>   shell Multi 기동시 구분 번호."
	echo "        [Json File]  Json config file (필수 입력 항목 아님)"
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
TMP_PATH=$PROJ_PATH/tmp

# Arguments
NODE_NO=$1
SEND_MAX=$2
SEQ_NO=$3
JSON_CONFIG_FILE=$4

# Global Var.
g_ret_isCreateBlock=""

JQ_BIN="$PROJ_PATH/bin/jq"

if [ -z $JSON_CONFIG_FILE ];
then
	JSON_CONFIG_FILE="sendTx-rpc.json"
fi

JSON_FILE="$SETUP_INFO_PATH/$JSON_CONFIG_FILE"

TIMESTAMP_FILE=$TLOG_PATH/timestamp_rpc.log
CMD_SAVE_FILE=$TMP_PATH/cmd-$NODE_NO-$SEQ_NO.tmp

TLOG_FILE_NAME="${TLOG_PATH}/${NODE_NO}-${SEQ_NO}-${BASE_NAME/\.sh/}.log"

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
RPC_PORT=$(( ${RPC_PORT} + ${NODE_NO#0} ))
HTTP_CONNECT="http://$RPC_IP:$RPC_PORT"
TX_TYPE=`func_JsonParsing ".arguments.tx_type"`

TX_FROM=`func_JsonParsing ".node_${NODE_NO#0}.${TX_TYPE}.from"`
TX_TO=`func_JsonParsing ".node_${NODE_NO#0}.${TX_TYPE}.to"`
TX_GAS=`func_JsonParsing ".node_${NODE_NO#0}.${TX_TYPE}.gas"`
TX_GASPRICE=`func_JsonParsing ".node_${NODE_NO#0}.${TX_TYPE}.gasPrice"`
TX_VALUE=`func_JsonParsing ".node_${NODE_NO#0}.${TX_TYPE}.value"`
TX_DATA=`func_JsonParsing ".node_${NODE_NO#0}.${TX_TYPE}.data"`

# Decimal -> Hex
TX_GAS=`func_Decimal2Hex $TX_GAS`
TX_GASPRICE=`func_Decimal2Hex $TX_GASPRICE`
TX_VALUE=`func_Decimal2Hex $TX_VALUE`



#####################################################
# Shared variable
#----------------------------------------------------
g_ret_getBalance=""
g_ret_last_sendTx_sub=""


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
# 잔고 조회 : geth
#----------------------------------------------------
func_getBalance_geth()
{
addr=$1
cmd="eth.getBalance(\"$addr\")"

result=`$BIN_PATH/geth attach rpc:$HTTP_CONNECT << EOF
$cmd
EOF`

g_ret_getBalance=`func_getValueParsing "$result"`
echo "$addr : $g_ret_getBalance"
}

#####################################################
# Tx 가 Blockchain 에 들어갔는지 확인
#----------------------------------------------------
func_isCreateBlock()
{
txhash=$1
cmd="eth.getTransactionReceipt(\"$txhash\")"

ret=`$BIN_PATH/geth attach rpc:http://$RPC_IP:$RPC_PORT << EOF
$cmd
EOF`

#echo "ret = $ret"

g_ret_isCreateBlock=""

while true
do
	for i in $ret ; do
		if [[ $i =~ "$txhash" ]];
		then
			g_ret_isCreateBlock="true"
			return
		elif [ "$i" == "null" ];
		then
			g_ret_isCreateBlock="false"
			return
		fi
	done
done

g_ret_isCreateBlock="false"
}

#####################################################
# Tx 전송 및 응답 처리
#----------------------------------------------------
func_last_sendTx_sub()
{
cmd="$1"

ret=`$BIN_PATH/geth attach rpc:http://$RPC_IP:$RPC_PORT << EOF
$cmd
EOF`

g_ret_last_sendTx_sub=""

for i in $ret ; do
	if [[ $i =~ "\"0x" ]];
	then
		g_ret_last_sendTx_sub="$i"
		return
	fi
done

echo "func_last_sendTx_sub : failed"
g_ret_last_sendTx_sub="null"
}

#####################################################
func_last_sendTx()
{
cmd=$1
while true
do
	func_last_sendTx_sub "$cmd"

	if [ "$g_ret_last_sendTx_sub" != "null" ];
	then
		break
	fi
done

txHash=$g_ret_last_sendTx_sub

retryCnt=0
while true
do
	retryCnt=`expr $retryCnt + 1`
	func_isCreateBlock "$txHash"

	if [ "$g_ret_isCreateBlock" == "true" ];
	then
		break
	fi

	if [ $retryCnt -ge 30 ];
	then
		break
	fi

	sleep 0.1

done
}

#####################################################
# GetTransactionReceipt processing
#----------------------------------------------------
func_receipt()
{
txhash_found_cnt=0

# file read
exec < $TLOG_FILE_NAME
while read line
do
	if [[ "$line" =~ "SEND_END" ]];
	then
		break
	fi

	if [[ "$line" =~ "\"0x" ]];
	then
		tx_hash=$line
		tx_hash=${tx_hash//\"/} # " 제거

		#echo "tx_hash = $tx_hash"

		while true
		do
			func_isCreateBlock $tx_hash

			if [ "$g_ret_isCreateBlock" == "true" ];
			then
				txhash_found_cnt=$(( $txhash_found_cnt + 1 ))
				echo "[$txhash_found_cnt] TxHash = $tx_hash"
				break
			fi

			#echo "Not found. TxHash($tx_hash)"
			sleep 1
		done
	fi

done

echo "---------------------------------------------------------------"
echo "Tx 완료 / 전체 : $txhash_found_cnt / $SEND_MAX"
echo "---------------------------------------------------------------"
}

#####################################################
# SendTransaction processing
#----------------------------------------------------
func_sendTx()
{
arg_loop_max=$1
arg_type=$2
arg_from=$3
arg_to=$4
arg_value=$5
arg_data=$6

echo "==============================================================="
func_getBalance_geth $arg_from
func_getBalance_geth $arg_to

cmd_one="personal.unlockAccount(\"$arg_from\", \"reapchain\", 0)"

if [ "$arg_type" == "messages" ];
then
	cmd_fmt="eth.sendTransaction({from: \"$arg_from\", to: \"$arg_to\", value: \"$arg_value\"})"
else
	cmd_fmt="eth.sendTransaction({from: \"$arg_from\", to: \"$arg_to\", value: \"$arg_value\", data: \"$arg_data\"})"
fi

# Tx 송신 & 수신(접수) 처리
echo "$cmd_one" > $CMD_SAVE_FILE

loop_cnt=0
while [ $loop_cnt -lt $arg_loop_max ];
do
	loop_cnt=$(( $loop_cnt + 1 ))

	echo "$cmd_fmt" >> $CMD_SAVE_FILE
done

cmd=$(<$CMD_SAVE_FILE)
rm -f $CMD_SAVE_FILE

beginTime=$(date +%s%N)

echo "S0-$NODE_NO-$SEQ_NO $(date +"%Y-%m-%d %H:%M:%S.%3N")" >> $TIMESTAMP_FILE

$BIN_PATH/geth attach rpc:$HTTP_CONNECT << EOF
$cmd
EOF

endTime1=$(date +%s%N)

echo "E1-$NODE_NO-$SEQ_NO $(date +"%Y-%m-%d %H:%M:%S.%3N")" >> $TIMESTAMP_FILE

echo "---------------------------------------------------------------"
elapsed=`echo "($endTime1 - $beginTime) / 1000000" | bc `
elapsedSec=`echo "scale=6;$elapsed / 1000" | bc | awk '{printf "%.6f", $1}'`
echo "Total elapsed : $elapsedSec sec"

echo "SEND_END"

#func_last_sendTx "$cmd_fmt"
func_receipt

endTime2=$(date +%s%N)
echo "E2-$NODE_NO-$SEQ_NO $(date +"%Y-%m-%d %H:%M:%S.%3N")" >> $TIMESTAMP_FILE

echo "---------------------------------------------------------------"
elapsed=`echo "($endTime2 - $beginTime) / 1000000" | bc `
elapsedSec1=`echo "scale=6;$elapsed / 1000" | bc | awk '{printf "%.6f", $1}'`

elapsed=`echo "($endTime2 - $endTime1) / 1000000" | bc `
elapsedSec2=`echo "scale=6;$elapsed / 1000" | bc | awk '{printf "%.6f", $1}'`

echo "Total elapsed : $elapsedSec1 sec ($elapsedSec2)"

echo "---------------------------------------------------------------"
func_getBalance_geth $arg_from
func_getBalance_geth $arg_to
}

#####################################################
# main
#####################################################

func_sendTx $SEND_MAX $TX_TYPE $TX_FROM $TX_TO $TX_VALUE $TX_DATA > $TLOG_FILE_NAME

