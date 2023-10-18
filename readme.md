# AllChat

Open source universal chat window.

## Quick Start

**Docker Compose**
```yaml
services:
  allchat:
    image: ghcr.io/graytonio/allchat:latest
    volumes:
      - ${PWD}/config.yaml:/config.yaml
    ports:
      - "8080:8080"
```

## Config File

The configuration file handles which chats to connect to and configuring them

```yaml
log_level: debug
enabled_chats:
  - twitch

twitch:
  channel: graytonio
```

| Key | Default | Description |
| -- | -- | -- |
| log_level | info | Changes the verbosity of loggin |
| enabled_chats | [] | Array of chat services to enable. Config for each of them is under the top level key of the same name |
| twitch.channel | "" | Name of twitch channel to pull chat from |

## Currently Supports:
    - [x] Twitch

## Future Updates:
    - [ ] YouTube

