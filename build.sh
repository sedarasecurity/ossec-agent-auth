#!/bin/bash

printf "** Building linux/amd64\n"
go build -a -o bin/linux-amd64/agent-auth github.com/sedarasecurity/ossec-agent-auth

printf "** Building linux/386\n"
go-linux-386 build -a -o bin/linux-386/agent-auth github.com/sedarasecurity/ossec-agent-auth

printf "** Building windows/386\n"
go-windows-386 build -a -o bin/windows-386/agent-auth.exe github.com/sedarasecurity/ossec-agent-auth

printf "** Building windows/amd64\n"
go-windows-amd64 build -a -o bin/windows-amd64/agent-auth.exe github.com/sedarasecurity/ossec-agent-auth
