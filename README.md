# Whatever

A simple implementation of Encrypted Tunnel

Extreme performance.

- [x] there's currently no UDP support yet

## Basic Usage

Decompress and execute the binary.

### Server

```whatever -s 'your-password@:9559'```

### Client

```whatever -c 'your-password@[server_address]:9559' -socks :1080```

## Advanced Usage

- [x] Replay Attack Mitigation
