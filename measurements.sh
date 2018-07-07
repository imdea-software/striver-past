#!/bin/bash

# killall striver &>/dev/null || true
PROG=$1
ARG2=$2
MAXEVS=$3
./striver $PROG $ARG2 $MAXEVS > /dev/null &
PID=$!

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

echo MEMORY $(($MEMORY / $COUNTER))
