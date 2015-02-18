#!/bin/bash

./build.sh

VERSION=$(git tag | tail -1)

tar -czf "_releases/sedara-agentauth-linux-amd64-${VERSION}.tar.gz" bin/linux-amd64/agent-auth
tar -czf "_releases/sedara-agentauth-linux-386-${VERSION}.tar.gz" bin/linux-386/agent-auth
7z a "_releases/sedara-agentauth-windows-386-${VERSION}.zip" bin/windows-386/agent-auth.exe
7z a "_releases/sedara-agentauth-windows-amd64-${VERSION}.zip" bin/windows-amd64/agent-auth.exe
