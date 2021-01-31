#!/bin/sh -eux

suffix="${GOARCH}-${GOOS}"

if [ "${GOOS}" = windows ]; then
    suffix="${suffix}.exe"
fi

agent_filename="rtun-${suffix}"
server_filename="rtun-server-${suffix}"

go build -o "${agent_filename}" ./agent/cmd
go build -o "${server_filename}" ./server/cmd

echo "::set-output name=agent::${agent_filename}"
echo "::set-output name=server::${server_filename}"
