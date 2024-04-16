#!/bin/bash
while true :
do
curl localhost:8088/random
#curl --location --request GET 'localhost:8080/training/v1/notebook/list' --header 'userID: wangbaoyi1' --header 'orgId: 010700' --header 'email: wangbaoyi1@jd.com' --header 'roleType: admin' --header 'req_seq_no: dxxcfdsf'
echo "success"
sleep 1
done
