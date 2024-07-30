# zenith
engineer's itch to build Redis client and server from scratch

## Provisions
- CLI to communicate to the kv database
- Redis like RESP specification to parse and execute commands
- Toggle to enable the kv storage to be volatile or persistent
- Utilise Redis based AOF(append only file) to record transactions as WALs
- Available commands as of first draft:
    - SET
    - GET
    - PING
    - ECHO
    - DEL
- More commands to be included post kv persistence 
- Render Redis CLI based errors via custom errors module
