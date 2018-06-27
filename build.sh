#!/bin/bash

pwd=`pwd`
GOPATH=$GOPATH:${pwd}

go build test_api.go
