#!/bin/bash -il
set -xe

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

version=$1

if [ -z $version ]; then
    version=v0.1
fi

go build -o crawl_photo main.go
chmod +x crawl_photo

tar -zcvf crawl_photo-linux-amd64.tar.gz \
  crawl_photo README.md
