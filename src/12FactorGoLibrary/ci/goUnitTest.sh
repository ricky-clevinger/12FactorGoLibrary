#!/bin/sh

set -e -u -x

export GOPATH=$PWD
export PATH=$PATH:$GOPATH

cd src/12FactorGoLibrary

go test