#!/bin/bash

set -e -u -x

export GOPATH=$PWD/src

go build -o built-project/library-go-app github.com/christopherneuhardt/12FactorGoLibrary/src