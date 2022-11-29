#!/bin/bash

expect -c "
        set timeout 200;
        spawn docker login --username=100009544039 ccr.ccs.tencentyun.com
        expect {
                \"Password:\" {send \"Cuczhangyi860830\n\";}
        }
expect eof;"

docker build -t ccr.ccs.tencentyun.com/dafangdocker/coros_strava .

docker push ccr.ccs.tencentyun.com/dafangdocker/coros_strava



#run 
# docker run -d  --restart=always --name coros_strava ccr.ccs.tencentyun.com/dafangdocker/coros_strava