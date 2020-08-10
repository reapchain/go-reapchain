qmanager 분리되었음.

make all 하면 build/bin/qman 생성됨

start-bootnode.sh 에

#!/bin/sh

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
BIN_PATH=$PROJ_PATH/bin
LOG_PATH=$PROJ_PATH/log

BOOTNODE_KEY="$SETUP_INFO_PATH/boot.key"
#BOOTNODE_KEY="$SETUP_INFO_PATH/bootnode.key"
BOOTNODE_PORT=30391
BOOTNODE_ADDR=192.168.0.80
QMANNODE_ADDR=192.168.0.80
QMANNODE_KEY="$SETUP_INFO_PATH/nodekey/nodekey.qman"  <= 기존 qman key를 그대로 사용
QMANPORT=30500

# bootnode
nohup $BIN_PATH/bootnode -nodekey $BOOTNODE_KEY -verbosity 9 -addr $BOOTNODE_ADDR:$BOOTNODE_PORT 2>&1 > $LOG_PATH/bootnode.log &
nohup $BIN_PATH/qman -qmankey $QMANNODE_KEY  -addr $QMANNODE_ADDR:$QMANPORT 2>&1 > $LOG_PATH/QManagerServer.log &   <= 추가함.


-------  기타 다른 쉘에세는 0부터 시작이 아니라, 1부터 시작하게끔 바꿈. --------------
## node startup
for ((i=1;i<$NODE_MAX;i++));
do
$PROJ_PATH/start-geth-node.sh $i
done