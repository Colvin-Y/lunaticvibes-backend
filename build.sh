#!/bin/bash
echo "已废弃，请移步 lunaticvibes-server!"
exit 0
rm -rf output
mkdir -p output
go build
mv lunaticvibes-backend output