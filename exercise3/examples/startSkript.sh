#!/bin/sh

IP="127.0.0.1"
MANAGER_A_PORT=15100
MANAGER_B_PORT=15101
MANAGER_A="127.0.0.1:15100"
MANAGER_B="127.0.0.1:15101"

FIRST_USER_PORT=15200

./fileManager -filename FileA.txt -force -ipAddress ${IP} -port ${MANAGER_A_PORT} > Manager1_out.txt
./fileManager -filename FileB.txt -force -ipAddress ${IP} -port ${MANAGER_B_PORT} > Manager2_out.txt

for i in 1 2; do
  CURRENT_PORT=$((${FIRST_USER_PORT}+1))
  ./fileUser --ipAddress 127.0.0.1 -port ${CURREN_TPORT} -managerA ${MANAGER_A} -managerB ${MANAGER_B} > fileUser${1}_out.txt
done
