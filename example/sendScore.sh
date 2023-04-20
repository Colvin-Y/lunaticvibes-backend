#!/bin/bash
curl -X POST \
	-F 'userID=1' \
	-F 'isCourse=0' \
	-F 'songHash=1ffb42c741af167beccd7b4352daec39' \
	-F 'courseHash=' \
    -F 'clearType=HARD_CLEAR' \
    -F 'lnMode=LN' \
    -F 'score=9999' \
    -F 'scoreMax=10000' \
    -F 'scoreRate=99.99' \
    -F 'scoreRank=AAA' \
    -F 'scorePG=4999' \
    -F 'scoreGR=1' \
    -F 'scoreGD=0' \
    -F 'scoreBD=0' \
    -F 'scorePR=0' \
    -F 'combo=5000' \
    -F 'laneOp=NORMAL' \
    -F 'gaugeOp=NORMAL' \
    -F 'inputType=CONTROLLER' \
    -F 'file=@/root/data/code/lunaticvibes-backend/example/test.txt' \
	http://localhost:8088/score