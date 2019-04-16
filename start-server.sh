#!/bin/bash

if [[ "${ROLE}" == "gateway" ]]
then
    /usr/bin/gateway
elif [[ "${ROLE}" == "rpcserv" ]]
then
    /usr/bin/rpcserv
else
    exit 1
fi
