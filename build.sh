#!/bin/bash

go mod tidy
rm -rf openai-on-wechat

GOOS=linux GOARCH=amd64 go build -o openai-on-wechat
