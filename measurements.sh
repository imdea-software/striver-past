#!/bin/bash

SEC=30

# killall striver-go &>/dev/null || true
timeout -s SIGINT ${SEC}s ./striver-go AVGK 10 > mes &
sleep 1
PID=`ps | grep striver | awk '{print $1}'`

# echo $PID

COUNTER=0
RES=0

# while [ $COUNTER -lt $SEC ];
while kill -0 $PID 2> /dev/null;
do
    COUNTER=$(($COUNTER+1))
    MEMORY=$(($MEMORY + `ps -p $PID -o vsize |grep -v VSZ`))
    sleep 1
done
PE=$(cat mes | grep Processed | awk '{print $3}')
rm mes

echo EVRATIO $(($PE / $SEC))
echo MEMORY $(($MEMORY / $COUNTER))
