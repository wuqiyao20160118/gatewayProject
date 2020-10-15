#!/bin/sh

export GO111MODULE=on
export GOPROXY=https://goproxy.io
go build -o bin/gatewayProject
ps aux | grep gatewayProject | grep -v 'grep' | awk '{print $2}' | xargs kill
nohup ./bin/gatewayProject -config=./conf/prod/ -endpoint=dashboard >> logs/dashboard.log 2>&1 &
echo 'nohup ./bin/gatewayProject -config=./conf/prod/ -endpoint=dashboard >> logs/dashboard.log 2>&1 &'
nohup ./bin/gatewayProject -config=./conf/prod/ -endpoint=server >> logs/server.log 2>&1 &
echo 'nohup ./bin/gatewayProject -config=./conf/prod/ -endpoint=server >> logs/server.log 2>&1 &'