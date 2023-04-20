#!/bin/bash
curl -X POST \
	-F 'userName=Darkzero' \
	-F 'userNickname=IceRabb' \
	-F 'userEmail=12315@qq.com' \
	-F 'userPwd=12315' \
	http://localhost:8088/signup