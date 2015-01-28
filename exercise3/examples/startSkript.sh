#!/bin/sh

IP="127.0.0.1"
MANAGER_A_PORT=15100
MANAGER_B_PORT=15101
MANAGER_A="127.0.0.1:15100"
MANAGER_B="127.0.0.1:15101"

FIRST_USER_PORT=15200

./fileManager
./fileManager

for i in 1 2 3 4 5 6 7 8 9; do
  ./fileUser
done
