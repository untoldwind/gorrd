#!/bin/bash

if [ ! -f "power.rrd" ];
then
    rrdtool create power.rrd \
       --start now-10s --step 1 \
       DS:watts:GAUGE:300:0:24000 \
       RRA:AVERAGE:0.5:1:864000 \
       RRA:AVERAGE:0.5:60:129600 \
       RRA:AVERAGE:0.5:3600:13392 \
       RRA:AVERAGE:0.5:86400:3660
fi

NOW=$(date +%s)

for i in $(seq 0 1 10000);
do
    let sec=$NOW+$i
    echo $i
    rrdtool update power.rrd $sec:$i
done
