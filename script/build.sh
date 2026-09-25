#!/bin/bash

if [[ -f ./h8 ]]; then
    rm ./h8
fi

go build -o h8 . && \
    ./h8 -i "./test/"

