#!/bin/bash
rm -rf output
mkdir -p output
go build
mv lunaticvibes-backend output