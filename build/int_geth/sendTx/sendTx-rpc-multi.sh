#!/bin/bash

func_Usage()
{
    echo ""
	echo "Usage : $(basename "$0") <Parallel Count> <Multi Count> <Send Count> [Json File]"
	echo "        <Parallel Count> Tx 전송 node 대상 (Min 2)"
	echo "                         30   ==> node 1 ~ 29(30-1)"
	echo "                         2:29 ==> node 2 ~ 29"
	echo "        <Multi Count>    1개 node 에 Tx 전송 Shell script 동시 기동 개수."
	echo "        <Send Count>     Tx 송신 횟수."
	echo "        [Json File]      Json config file."
    echo ""
}

if [ $# -lt 3 ]; then
	func_Usage
    exit 0
fi

# Parallel Count 파싱
PARALLEL_COUNT=$1

if [[ "$1" =~ ":" ]]; then
	FIRST_NODE_NO=$(echo "$PARALLEL_COUNT" | awk -F ":"  '{printf $1}')
	LAST_NODE_NO=$(echo "$PARALLEL_COUNT" | awk -F ":"  '{printf $2}')
else
	if [ $PARALLEL_COUNT -lt 2 ]; then
		func_Usage
		echo "\n<Parallel Count> 는 2보다 작을 수 없습니다.\n"
		exit 0
	fi
	FIRST_NODE_NO=1
	LAST_NODE_NO=$(( $1 - 1 ))
fi

if [ $LAST_NODE_NO -lt $FIRST_NODE_NO ]; then
	func_Usage
	echo "\n<Parallel Count> : LAST($LAST_NODE_NO)는 FIRST($FIRST_NODE_NOi) 보다 작을 수 없습니다.\n"
	exit 0
fi

MULTI_COUNT=$2
SEND_COUNT=$3
JSON_CONFIG_FILE=$4

BASE_NAME=$(basename "$0")
DIR_NAME=$(dirname "$0")
PROJ_PATH="$( cd "$DIR_NAME" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
TLOG_PATH=$PROJ_PATH/tlog
TMP_PATH=$PROJ_PATH/tmp

JQ_BIN="$PROJ_PATH/bin/jq"

if [ "$JSON_CONFIG_FILE" == "" ];
then
	JSON_CONFIG_FILE="sendTx-rpc.json"
fi
JSON_FILE="$SETUP_INFO_PATH/$JSON_CONFIG_FILE"

#####################################################
# JSON Parsing
#----------------------------------------------------
func_JsonParsing()
{
	ret_val=`cat $JSON_FILE | $JQ_BIN $1`
	ret_val=${ret_val//\"/} # 앞뒤의 " 제거
	echo $ret_val
}

CALL_SHELL=`func_JsonParsing ".shell_script"`
RPC_IP=`func_JsonParsing ".arguments.ip"`
RPC_PORT=`func_JsonParsing ".arguments.rpc_port"`

# RPC_IP 와 RPC_PORT 는 수행전 확인용
HTTP_CONNECT="http://$RPC_IP:$RPC_PORT"

TLOG_FILE_NAME="${TLOG_PATH}/0-0-${BASE_NAME/\.sh/}.log"

#####################################################
# PATH 체크 및 생성
#----------------------------------------------------
func_path_checking()
{
	if [ ! -d $TLOG_PATH ];
	then
		mkdir $TLOG_PATH
	fi

	if [ ! -d $TMP_PATH ];
	then
		mkdir $TMP_PATH
	fi
}

#####################################################
func_path_init()
{
	PATH_CNT=`ls $TLOG_PATH | wc -l`
	if [ $PATH_CNT -gt 0 ];
	then
		rm -f ${TLOG_PATH}/*
	fi

	PATH_CNT=`ls $TLOG_PATH | wc -l`
	if [ $PATH_CNT -gt 0 ];
	then
		rm -f ${TMP_PATH}/*
	fi
}

#####################################################
# main processing
#----------------------------------------------------
func_main()
{
	echo "CALL_SHELL       = $CALL_SHELL"
	if [ $FIRST_NODE_NO -eq $LAST_NODE_NO ]; then
		echo "PARALLEL_COUNT   = $FIRST_NODE_NO"
	else
		echo "PARALLEL_COUNT   = $FIRST_NODE_NO ~ $LAST_NODE_NO"
	fi
	echo "MULTI_COUNT      = $MULTI_COUNT"
	echo "SEND_COUNT       = $SEND_COUNT"
	echo "JSON_CONFIG_FILE = $JSON_CONFIG_FILE"

	node_nos=`seq -w $FIRST_NODE_NO $LAST_NODE_NO`
	seq_nos=`seq -w 1 $MULTI_COUNT`
	#echo "seq_nos = $seq_nos"

	for j in $node_nos; do
		for i in $seq_nos; do
			echo "$PROJ_PATH/$CALL_SHELL $j $SEND_COUNT $i $JSON_CONFIG_FILE"
			nohup $PROJ_PATH/$CALL_SHELL $j $SEND_COUNT $i $JSON_CONFIG_FILE &
		done
	done
}

#####################################################
# main
#####################################################

echo ""
read -n 1 -p "$HTTP_CONNECT 전송 목적지가 맞습니까? y/N > "
echo ""
[[ $REPLY != [yY] ]] && exit 0

func_path_checking
func_path_init

func_main > $TLOG_FILE_NAME

