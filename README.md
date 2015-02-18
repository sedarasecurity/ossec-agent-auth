ossec-agent-auth
================

Implementation of the agent-auth client in go

## Download:
[Latest Release](https://github.com/sedarasecurity/ossec-agent-auth/releases/latest)

## Usage:

```
Usage of agent-auth:
  -config="": Path to OSSEC config file (ossec.conf)
  -keyfile="": Path to OSSEC client keys file (client.keys)
  -listen=false: Enables running in server mode
  -manager="": Manager IP Address
  -name="localhost": Agent name
  -port=1515: Manager port
```

## Implementation:
By default, the agent will modify the ossec.conf file to update the manager's IP with the one given on the command line.

## Example:

```
$> agent-auth -config="/var/ossec/etc/ossec.conf" -keyfile="/var/ossec/etc/client.keys" -manager="192.168.0.2" -name="server1"
```

## Known Issues:
* DNS resolution currently not implemented; the manager flag expects an IP

## TODO:
* Certificate verification
* DNS resolution
* Configuration backups
