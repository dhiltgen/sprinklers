#!/bin/bash

# Crude curl script to fire up the sprinklers for the lawn

HEAD="Content-Type: application/json"

URL=http://sprinklers:80/circuits

curl -H "${HEAD}" -X POST -d '{"gpio": 22, "state":true, "remaining":"15m"}' ${URL}
sleep 15m
curl -H "${HEAD}" -X POST -d '{"gpio": 8, "state":true, "remaining":"15m"}' ${URL}
sleep 15m
curl -H "${HEAD}" -X POST -d '{"gpio": 9, "state":true, "remaining":"15m"}' ${URL}
sleep 15m
curl -H "${HEAD}" -X POST -d '{"gpio": 24, "state":true, "remaining":"15m"}' ${URL}
