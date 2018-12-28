#!/bin/bash

VERSION=0.0.2

rm -rf rootfs
docker build -t docker-redis-log-driver .
ID=$(docker create docker-redis-log-driver true)
mkdir rootfs
docker export $ID | tar -x -C rootfs/
docker plugin disable redis-log-driver:$VERSION
docker plugin rm redis-log-driver:$VERSION
docker plugin create redis-log-driver:$VERSION .
docker plugin enable redis-log-driver:$VERSION
