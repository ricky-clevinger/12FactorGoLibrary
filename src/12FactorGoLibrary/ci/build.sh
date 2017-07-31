#!/bin/sh
set -e -u -x
export GOPATH=$PWD
export PATH=$PATH:$GOPATH
#export LIBRARY="cgidevlib:Password1@tcp(cgiprojdevlibrary.cxyeb3wmov3g.us-east-1.rds.amazonaws.com:9871)/cgiprojdevlibrary"
cd src
cd 12FactorGoLibrary
go build index.go
