#!/bin/bash
BINARY_PATH=/home/go-binary
BINARY_NAME=cquec-roller

cd $BINARY_PATH
PID=$(pgrep -f ${BINARY_NAME})

if [ -z $PID ]; then
        echo "프로세스가 실행되고 있지 않습니다."
else
        echo "대상 프로세스 $PID를 Kill 처리하였습니다."
        kill -9 $PID
        sleep 5
fi

rm cquec-roller
mv cquec-roller-build cquec-roller

chmod +x ./$BINARY_NAME
nohup ./$BINARY_NAME > nohup.out 2>&1 &
