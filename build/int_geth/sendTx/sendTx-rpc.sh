#!/bin/sh

if [ $# -ne 1 ]; then
	echo ""
	echo "Usage : $(basename "$0") <Node No>"
	echo ""

	return
fi

#MY_LOCAL_IP=$(ifconfig enp0s3 | grep inet | grep netmask | awk {'print $2'})
MY_LOCAL_IP=$(hostname -I | awk '{print $1}')

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
BASE_NAME="$(basename "$0")"
BIN_PATH=$PROJ_PATH/bin
LOG_PATH=$PROJ_PATH/log
NODE_NAME=node$1
RPC_PORT=`expr 8540 \+ "$1"`
RPC_IP=localhost
#RPC_IP=$MY_LOCAL_IP
#RPC_IP=192.168.0.86
#RPC_IP=192.168.0.100


while true
do
$BIN_PATH/geth attach rpc:http://$RPC_IP:$RPC_PORT << EOF
personal.unlockAccount(eth.accounts[0], "reapchain")
eth.sendTransaction({from: eth.accounts[0], to: eth.accounts[1], value: web3.toWei(1, "ether")})
EOF
sleep 1
done >> $LOG_PATH/$BASE_NAME.log

