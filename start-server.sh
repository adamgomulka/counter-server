#!/bin/bash

if [[ "${ROLE}" == "gateway" ]]
then
    /usr/bin/gateway
elif [[ "${ROLE}" == "rpcserv" ]]
then
    /usr/bin/rpcserv
elif [[ "${ROLE}" == "tcp_wait" ]]
then
    /usr/bin/tcp_wait ${WAIT_SERVICE} ${WAIT_PORT}
else
    exit 1
fi
