#!/usr/bin/env bash

pkill bimgserver-amd64
nohup ./bimgserver-amd64 > bimgserver.log 2>&1 &