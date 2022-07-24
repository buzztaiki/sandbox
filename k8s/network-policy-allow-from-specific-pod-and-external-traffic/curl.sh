#!/bin/bash -x

kubectl exec "$1" -- curl -s -I --connect-timeout 2 http://nginx:8000
