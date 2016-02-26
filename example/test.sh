#!/bin/sh

for i in `seq 1 50`; do
        curl -X POST http://127.0.0.1:8000/work --data "delay=5s&name=Test"
done

for i in `seq 1 50`; do
        curl -X GET http://127.0.0.1:8000/hello
done
