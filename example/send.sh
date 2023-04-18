#!/bin/bash
curl -X POST \
	-F 'score=90' \
	-F 'userID=123' \
	-F 'file=@/root/data/code/lunaticvibes-backend/example/test.txt' \
	http://localhost:8088/score
