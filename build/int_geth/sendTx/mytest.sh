#!/bin/bash

if [ $# -ne 2 ]; then
        echo ""
        echo "Usage : $(basename "$0") <Node No> <Loop count>"
        echo ""

        return
fi

#MY_LOCAL_IP=$(ifconfig enp0s3 | grep inet | grep netmask | awk {'print $2'})

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
BASE_NAME="$(basename "$0")"
BIN_PATH=$PROJ_PATH/bin
LOG_PATH=$PROJ_PATH/log
NODE_NAME=node$1
RPC_PORT=`expr 8540 \+ "$1"`
#RPC_IP=localhost
RPC_IP=192.168.0.86

loop_max=$2
loop_cnt=0

echo "personal.unlockAccount(eth.accounts[0], \"reapchain\", 0)" >> command.txt

while [ $loop_cnt -lt $loop_max ];
do
loop_cnt=$(( $loop_cnt + 1 ))
echo "eth.sendTransaction({from: eth.accounts[0], to: eth.accounts[1], value: web3.toWei(1, \"reap\")})" >> command.txt
done

echo "" >> command.txt

command=$(<command.txt)
echo "command = $command"

$BIN_PATH/geth attach rpc:http://$RPC_IP:$RPC_PORT << EOF > $LOG_PATH/$BASE_NAME.log
$command
EOF

rm -f command.txt

