#!/bin/bash
PROJECT="route"
BIN_FILE="route-registrar"
rm -rf ${PROJECT}
rm ${PROJECT}.zip
mkdir ${PROJECT}
rm ${BIN_FILE}
go build
if [ $? -eq 0 ];then
    cp manifest.yml ${PROJECT}
    cp ${BIN_FILE} ${PROJECT}
    cp -rf conf/ ${PROJECT}
    cp -rf static/ ${PROJECT}
    cp -rf views/ ${PROJECT}
    cd ${PROJECT}
    zip -r ${PROJECT}.zip *
    mv  ${PROJECT}.zip ..
else
    echo "build failed"
    exit 1
fi
