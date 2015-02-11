#!/bin/bash

./build.sh

tar -czf releases/sedara-agentauth-linux-amd64.tar.gz bin/linux-amd64/agent-auth
tar -czf releases/sedara-agentauth-linux-386.tar.gz bin/linux-386/agent-auth
7z a releases/sedara-agentauth-windows-386.zip bin/windows-386/agent-auth.exe
7z a releases/sedara-agentauth-windows-amd64.zip bin/windows-amd64/agent-auth.exe
