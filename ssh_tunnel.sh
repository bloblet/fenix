#!/usr/bin/env bash

ssh -N -L 9042:localhost:9042 "$1"@vps.bloblet.com