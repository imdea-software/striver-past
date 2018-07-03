#!/bin/bash

SEC=60

timeout -s SIGINT ${SEC}s ./striver-go > mes &

PID=$!

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
