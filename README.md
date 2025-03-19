# Redis

## ⚠️ Warning

**This implementation is for educational purposes only and should NOT be used in production environments.** It does not implement all Redis security features and has not undergone security auditing or performance optimization for production workloads.

## Overview

This project is an implementation of a Redis server in Go, designed to help understand the inner workings of Redis, TCP networking, and in-memory databases. It follows the Redis Serialization Protocol (RESP) and implements core Redis functionality.

## Features

The implementation follows a phased approach:

- [x] Basic TCP server with PING command
- [ ] RESP protocol implementation
- [ ] Basic string operations (GET, SET)
- [ ] Data structures (Lists, Sets, etc.)
- [ ] Key expiration (TTL)
- [ ] Persistence (RDB snapshots, AOF logs)

## Architecture

The system is designed with a clean separation of concerns, divided into these primary components:

1. **Server**: Handles TCP connections with clients
2. **Protocol**: Manages serialization/deserialization of the RESP format
3. **Command Layer**: Routes commands to appropriate handlers
4. **Store**: Manages in-memory data structures and persistence


## Getting Started

### Prerequisites

- Go 1.24 or higher

### Installation

```bash
git https://github.com/benidevo/redis
cd redis
```

### Building

```bash
go build -o redis ./cmd/redis
```

Or use the included script:

```bash
sh scripts/run.sh
```

### Running

```bash
./redis
```

The server will start and listen on the default Redis port (6379).

## Implementation Details

The project implements Redis from first principles:

1. **TCP Communication**: Uses Go's `net` package for accepting and handling connections
2. **RESP Protocol**: Custom implementation of Redis protocol for command parsing
3. **Command Execution**: Dispatches commands to appropriate handlers
4. **Data Storage**: Implements in-memory data structures with optional persistence
