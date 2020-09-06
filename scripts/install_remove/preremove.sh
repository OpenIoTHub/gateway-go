#!/usr/bin/env bash

echo preremove.sh:
systemctl stop gateway-go
systemctl disable gateway-go