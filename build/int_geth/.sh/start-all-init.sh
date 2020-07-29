#!/bin/bash

if [ $# -ne 1 ]; then
    echo ""
    echo "Usage : $(basename "$0") <Node Max>"
    echo ""

    exit 0
fi

PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
SETUP_INFO_PATH=$PROJ_PATH/setup_info
NODEKEY_PATH=$SETUP_INFO_PATH/nodekey
NODE_MAX=$1

$PROJ_PATH/mykill
sleep 2

$PROJ_PATH/start-init-node.sh $NODE_MAX
sleep 1

$PROJ_PATH/create_address_add_all.sh $NODE_MAX

#$PROJ_PATH/import_address.sh 5 $SETUP_INFO_PATH/nodekey.user1
#$PROJ_PATH/import_address.sh 5 $SETUP_INFO_PATH/nodekey.user2

$PROJ_PATH/import_address.sh 6 $SETUP_INFO_PATH/nodekey.fee
$PROJ_PATH/import_address.sh 6 $SETUP_INFO_PATH/nodekey.governance
$PROJ_PATH/governancekey_copy.sh 6

$PROJ_PATH/import_address.sh 7 $SETUP_INFO_PATH/nodekey.coin_issuance
$PROJ_PATH/import_address.sh 7 $SETUP_INFO_PATH/nodekey.foundation
$PROJ_PATH/import_address.sh 7 $SETUP_INFO_PATH/nodekey.presale
$PROJ_PATH/import_address.sh 7 $SETUP_INFO_PATH/nodekey.alloc
$PROJ_PATH/import_address.sh 7 $SETUP_INFO_PATH/nodekey.contract_fee
$PROJ_PATH/import_address.sh 7 $SETUP_INFO_PATH/nodekey.alloc_incentive
$PROJ_PATH/import_address.sh 7 $SETUP_INFO_PATH/nodekey.alloc_send
$PROJ_PATH/import_address.sh 7 $SETUP_INFO_PATH/nodekey.batch_operation
sleep 1

$PROJ_PATH/start-geth-all.sh $NODE_MAX

