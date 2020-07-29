#!/bin/bash

if [ $# -ne 2 ]; then
	echo ""
	echo "Usage : $(basename "$0") <Node No> <Send Count>"
	echo "        <Send Count> When set to 0, infinite repeat."
	echo ""

	exit
fi

#MY_LOCAL_IP=$(ifconfig enp0s3 | grep inet | grep netmask | awk {'print $2'})
MY_LOCAL_IP=$(hostname -I | awk '{print $1}')
PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
BASE_NAME="$(basename "$0")"
BIN_PATH=$PROJ_PATH/bin
LOG_PATH=$PROJ_PATH/log
TLOG_PATH=$PROJ_PATH/tlog
NODE_NAME=node$1
RPC_PORT=$((  8540 + $1 ))
RPC_IP=$MY_LOCAL_IP
#RPC_IP=localhost
HTTP_CONNECT="http://$RPC_IP:$RPC_PORT"
TLOG_FILE_NAME="${TLOG_PATH}/${RPC_PORT}_${BASE_NAME}.log"


ADDR_IX0=0
ADDR_IX1=1
ADDR_IX2=2
SNED_BAL=1000000000000000000
SEND_MAX=$2

#####################################################
# Shared variable
g_ret_sendTx=""
g_ret_isCreateBlock=""
g_ret_getBalance=""
beginTime=""


#####################################################
# tlog 폴더 체크 및 생성
func_log_checking()
{
if [ ! -d $TLOG_PATH ];
then
	mkdir $TLOG_PATH
fi

}

#####################################################
func_getValueParsing()
{
arg=$1
findok=false
for i in $arg ; do
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
func_getBalance()
{
addr_index=$1
cmd="eth.getBalance(eth.accounts[$addr_index])"
#echo "getBalance : cmd = $cmd"
arg=`$BIN_PATH/geth attach rpc:http://$RPC_IP:$RPC_PORT << EOF
$cmd
EOF`
g_ret_getBalance=`func_getValueParsing "$arg"`
echo "eth.accounts[$addr_index] : $g_ret_getBalance"
}

#####################################################
func_isCreateBlock()
{
cmd="eth.getTransactionReceipt($1)"
#echo "func_isCreateBlock : cmd = $cmd"

ret=`$BIN_PATH/geth attach rpc:http://$RPC_IP:$RPC_PORT << EOF
$cmd
EOF`

g_ret_isCreateBlock=""

while true
do
	for i in $ret ; do
		#echo "func_isCreateBlock i = $i"
		if [ "$i" == "{" ];
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
func_sendTx()
{
cmd1="personal.unlockAccount(eth.accounts[$1], \"reapchain\")"
cmd2="eth.sendTransaction({from: eth.accounts[$1], to: eth.accounts[$2], value: web3.toWei($3, \"reap\")})"
#echo "func_sendTx : cmd1 = $cmd1"
echo "func_sendTx : cmd2 = $cmd2"

ret=`$BIN_PATH/geth attach rpc:http://$RPC_IP:$RPC_PORT << EOF
$cmd1
$cmd2
EOF`

g_ret_sendTx=""

for i in $ret ; do
	#echo "i = $i"
	if [[ $i =~ "\"0x" ]];
	then
		g_ret_sendTx="$i"
		return
	fi
done

echo "func_sendTx : failed"
g_ret_sendTx="null"
}

#####################################################
main_loop()
{
loop_arg=$1
loop_max=$1

if [ $loop_arg -eq 0 ];
then
	loop_max=1
fi

echo "==============================================================="
func_getBalance $ADDR_IX0
func_getBalance $ADDR_IX1
#func_getBalance $ADDR_IX2

loop_cnt=0
while [ $loop_cnt -lt $loop_max ];
do
	loop_cnt=$(( $loop_cnt + 1 ))

	if [ $loop_arg -eq 0 ];
	then
		loop_max=$(( $loop_cnt + 1 ))
	fi

	TIME_DISP=`date +"%Y-%m-%d %H:%M:%S.%3N"`
	echo "---------------------------------------------------------------"
	echo "[$loop_cnt] $TIME_DISP"

	beginTime=$(date +%s%N)

	while true
	do
		func_sendTx	$ADDR_IX0 $ADDR_IX1 $SNED_BAL

		if [ "$g_ret_sendTx" != "null" ];
		then
			break
		fi
	done

	echo "txHash = $g_ret_sendTx"

	txHash=$g_ret_sendTx

	retryCnt=0
	while true
	do
		retryCnt=`expr $retryCnt + 1`
		func_isCreateBlock "$txHash"

		if [ "$g_ret_isCreateBlock" == "true" ];
		then
			func_getBalance $ADDR_IX0

			func_getBalance $ADDR_IX1

			#func_getBalance $ADDR_IX2

			endTime=$(date +%s%N)
			elapsed=`echo "($endTime - $beginTime) / 1000000" | bc `
			elapsedSec=`echo "scale=6;$elapsed / 1000" | bc | awk '{printf "%.6f", $1}'`
			echo "Total elapsed : $elapsedSec sec"
			break
		fi

		if [ $retryCnt -ge 30 ];
		then
			break
		fi

		sleep 0.1

	done
done
}

#####################################################
## main
#####################################################

echo ""
read -n 1 -p "$HTTP_CONNECT 전송 목적지가 맞습니까? y/N > "
echo ""
[[ $REPLY != [yY] ]] && exit 0

func_log_checking

main_loop $SEND_MAX >> $TLOG_FILE_NAME

