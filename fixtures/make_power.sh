#!/bin/bash

if [ ! -f "power.rrd" ];
then
    rrdtool create power.rrd \
       --start now-10s --step 1s \
       DS:watts:GAUGE:5m:0:100000 \
       DS:counts:COUNTER:5m:0:100000 \
       RRA:AVERAGE:0.5:1s:12h \
       RRA:AVERAGE:0.5:1m:2d \
       RRA:AVERAGE:0.5:1h:20d \
       RRA:AVERAGE:0.5:1d:1M
fi

NOW=$(date +%s)

for i in $(seq 0 1 10000);
do
    let sec=$NOW+$i
    echo $i
    rrdtool update power.rrd $sec:$i:$i
done

for i in $(seq 10000 1 20000);
do
    let sec=$NOW+$i*5
    echo $i
    rrdtool update power.rrd $sec:$i:$i
done
