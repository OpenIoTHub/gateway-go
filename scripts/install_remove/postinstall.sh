#!/usr/bin/env bash

echo postinstall.sh
systemctl enable gateway-go
systemctl start gateway-go