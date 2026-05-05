#!/usr/bin/env bash

polo() {
    ./polo
    stty sane
    if [ -f "/tmp/goline-polo-path.txt" ]; then
        cd $(cat /tmp/goline-polo-path.txt)
        rm /tmp/goline-polo-path.txt
    fi
}

polo