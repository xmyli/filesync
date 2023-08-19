# File Synchronization Tool [WIP]

Two-way continuous file synchronization tool with a client-server architecture.

### Features
- File system notifications and periodic polling to detect changes.
- Delta encoding algorithm to save bandwidth by only transmitting the differences between files.
- Server-sent events to notify clients of changes on the server.

### Usage
1. Build cmd/server and cmd/client
    ```
    go build ./cmd/server
    go build ./cmd/client
    ```
2. Run the server executable, specifying a root directory
    ```
    ./server -d=FILES_DIRECTORY
    ```
3. Run the client executable, specifying a root directory
    ```
    ./client -d=FILES_DIRECTORY
    ```

### Todos
- [ ] Fix bugs
