#!/bin/sh

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

go build -o lambda main.go
if [ $? -ne 0 ]; then
    exit 1
fi

if [ -f lambda.zip ]; then
    rm lambda.zip
fi

zip -r lambda.zip ./lambda
if [ $? -ne 0 ]; then
    exit 1
fi