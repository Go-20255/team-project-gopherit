#!/bin/bash

polo() {
    ./polo
    stty sane
    if [ -f "/tmp/polo.txt" ]; then
        cd $(cat /tmp/polo.txt)
        rm /tmp/polo.txt
    fi
}

polo