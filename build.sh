#!/usr/bin/env bash 
# install packages and dependencies
go get github.com/gorilla/mux
go get github.com/aws/aws-sdk-go
go get html/template
go get encoding/json
go get net/http

# build command
go build -o bin/aws aws.go