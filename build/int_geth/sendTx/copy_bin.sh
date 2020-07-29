#!/bin/bash

if [ $# -ne 1 ]; then
	echo ""
	echo "Usage : $(basename "$0") <Binary Name>"
	echo ""

	exit 0
fi

# PROJ_PATH 는 shell script 가 위치한 PTAH 입니다.
PROJ_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
YYYYMMDD="$(date +'%Y%m%d-%H%M')"
BIN_PATH=$PROJ_PATH/bin
BIN_NAME=$1
TARGET_NAME="$1.$YYYYMMDD"

mv $BIN_PATH/$BIN_NAME $BIN_PATH/$TARGET_NAME
cp -p $GETH_BIN/$BIN_NAME $BIN_PATH

