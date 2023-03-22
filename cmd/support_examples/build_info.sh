#!/bin/bash

rm build.go

VERSION='1.2.3_123'

echo 'package main' > build.go
echo ''  >> build.go
echo "var version=\"${VERSION}\""  >> build.go

